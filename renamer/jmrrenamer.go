package renamer

import (
	"errors"
	"fmt"
	"log"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	m "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type JmrRenamer struct {
	mip m.MediaInfoProvider
	fs  filesystem.FileSystemProvider
	cfg config.ConfigProvider
}

func NewJmrRenamerV1(mip m.MediaInfoProvider, fs filesystem.FileSystemProvider, cfg config.ConfigProvider) *JmrRenamer {
	return &JmrRenamer{
		mip,
		fs,
		cfg,
	}
}

func (j *JmrRenamer) GetMediaNameAndYear(rawFilename string) MediaNameAndYear {
	// Step 1: remove special characters from input
	onlyBasename := path.Base(rawFilename)
	outputFilename := removeSpecialCharacters(onlyBasename)

	// Step 2: extract year from input, if present otherwise 0
	// also remove all string after the year (including year)
	outputYear := extractYear(&outputFilename)

	// Step 3: remove all double spaces due to previous steps
	outputFilename = removeDoubleWhitespace(outputFilename)

	return MediaNameAndYear{outputFilename, outputYear}
}
func (j *JmrRenamer) GetMediaSeasonAndEpisode(rawFilepath string) MediaSeasonAndEpisode {
	filepathWithoutExtension, _ := strings.CutSuffix(rawFilepath, path.Ext(rawFilepath))
	in := strings.ToLower(filepathWithoutExtension)
	in += "_" // added to pass the last two regexps

	testREs := []*regexp.Regexp{
		regexp.MustCompile(`s(\d{2})e(\d{2})`),                     // SXXEXX
		regexp.MustCompile(`season[ \-_]+(\d{1,2})[ \-_]+(\d{2})`), // Season X - XX
		regexp.MustCompile(`s(\d+)[ \-_]+(\d+)`),                   // SX - XX
		regexp.MustCompile(`(\d{1,2})x(\d{2})`),                    // [SxEE] and [SSxEE]
		regexp.MustCompile(`episode[ \-_]+(\d+)`),                  // Episode XX (No season information)
		regexp.MustCompile(`[^0-9](\d{1,2})[^0-9]`),                // XX (No season information)
	}

	for _, testre := range testREs {
		m1 := testre.FindStringSubmatch(in)
		if m1 != nil {
			if len(m1) == 2 {
				// contains only episode number
				episode, e2 := strconv.Atoi(m1[1])
				if e2 != nil {
					panic("Cannot parse episode number. " + e2.Error())
				}
				if episode > 200 {
					return MediaSeasonAndEpisode{-1, -1}
				}
				return MediaSeasonAndEpisode{1, episode}
			} else if len(m1) == 3 {
				// contains both episode and season number
				season, e1 := strconv.Atoi(m1[1])
				if e1 != nil {
					panic("Cannot parse season number. " + e1.Error())
				}
				episode, e2 := strconv.Atoi(m1[2])
				if e2 != nil {
					panic("Cannot parse episode number. " + e2.Error())
				}
				return MediaSeasonAndEpisode{season, episode}
			}
		}
	}

	return MediaSeasonAndEpisode{-1, -1}
}

func (j *JmrRenamer) SelectEntriesForRename(rootEntry filesystem.DirEntry, mediaType m.MediaType) EntriesAndIgnores {
	switch mediaType {
	case m.MediaTypeMovie:
		return j.selectEntriesForMovieRename(rootEntry)
	case m.MediaTypeTV:
		return j.selectEntriesForTVRename(rootEntry)
	default:
		panic(fmt.Sprintf("SelectEntriesForRename not implemented for mediaType: %s", mediaType))
	}
}

func (j *JmrRenamer) ConfirmEntriesForRename(entries RenameMediaConfirmRequest) (int, error) {
	allPathPairs := []filesystem.PathPair{}

	for _, entry := range entries {
		oldRoot := path.Dir(entry.Entry.Path)
		newRoot := path.Join(oldRoot, ".jmr-renames")
		newEntryDir := path.Join(newRoot, j.mip.GetJellyfinCompatibleDirectoryName(entry.Info))

		// Create newEntryDir
		if j.fs.CreateDirectory(newEntryDir) != true {
			// Failed to create newEntryDir
			return -1, errors.New("Cannot create entry directory " + newEntryDir)
		}

		for _, renEntry := range entry.Selected {
			newPath := j.renameSingleEntry(entry, renEntry, newEntryDir, false)
			allPathPairs = append(allPathPairs, filesystem.PathPair{
				OldPath: renEntry.Media.Path,
				NewPath: newPath,
			})

			if renEntry.Subtitle != nil {
				newPath := j.renameSingleEntry(entry, renEntry, newEntryDir, true)
				allPathPairs = append(allPathPairs, filesystem.PathPair{
					OldPath: renEntry.Subtitle.Path,
					NewPath: newPath,
				})
			}
		}
	}

	progress := make(chan []filesystem.FileTransferProgress)
	go j.fs.MoveFiles(allPathPairs, progress)

	for p := range progress {
		log.Println("Overall progress:")

		for _, pp := range p {
			log.Println(pp.ToString())
		}
		log.Println()
	}

	// Delete original source entries to save space and reduce duplication
	// TODO: uncomment this block before deploy
	// for _, entry := range entries {
	// 	if j.fs.DeleteDirectory(entry.Entry.Path) != true {
	// 		log.Printf("Cannot delete directory / file %s", entry.Entry.Path)
	// 	}
	// }

	return 1, nil
}

func (j *JmrRenamer) renameSingleEntry(entry RenameMediaResponseItem, e RenameEntry, newEntryDir string, isSubtitle bool) string {
	var ext string
	if isSubtitle {
		ext = path.Ext(e.Subtitle.Path)
	} else {
		ext = path.Ext(e.Media.Path)
	}

	if e.Season > 0 && e.Episode > 0 {
		// TV Show
		newBase := fmt.Sprintf("Season %02d/%s S%02dE%02d%s", e.Season, entry.Info.Name, e.Season, e.Episode, ext)
		return path.Join(newEntryDir, newBase)
	} else {
		// Movie
		newBase := fmt.Sprintf("%s%s", j.mip.GetJellyfinCompatibleDirectoryName(entry.Info), ext)
		return path.Join(newEntryDir, newBase)
	}
}

func (j *JmrRenamer) getTotalSizeForRenameEntries(entries RenameMediaConfirmRequest) int64 {
	// compute total size of all files
	totalSize := int64(0)
	for _, entry := range entries {
		for _, renEntry := range entry.Selected {
			totalSize += renEntry.Media.Size
			if renEntry.Subtitle != nil {
				totalSize += renEntry.Subtitle.Size
			}
		}
	}

	return totalSize
}

func (j *JmrRenamer) selectEntriesForMovieRename(rootEntry filesystem.DirEntry) EntriesAndIgnores {
	allMediaFiles := j.getAllEntriesOfExtension(rootEntry, j.cfg.GetMediaExtensions())
	slices.SortFunc(allMediaFiles, util.SortBySizeDesc)

	allSubtitleFiles := j.getAllEntriesOfExtension(rootEntry, j.cfg.GetSubtitleExtensions())
	slices.SortFunc(allSubtitleFiles, util.SortBySizeDesc)

	out := EntriesAndIgnores{
		Selected: []RenameEntry{},
		Ignored:  []filesystem.DirEntry{},
	}

	for i, mediaFile := range allMediaFiles {
		if i == 0 {
			out.Selected = append(out.Selected, RenameEntry{
				Media:    &mediaFile,
				Subtitle: nil,
			})
		} else {
			out.Ignored = append(out.Ignored, mediaFile)
		}
	}

	for i, srtFile := range allSubtitleFiles {
		if i == 0 && len(out.Selected) > 0 {
			out.Selected[0].Subtitle = &srtFile
		} else {
			out.Ignored = append(out.Ignored, srtFile)
		}
	}

	return out
}
func (j *JmrRenamer) selectEntriesForTVRename(rootEntry filesystem.DirEntry) EntriesAndIgnores {
	allMediaFiles := j.getAllEntriesOfExtension(rootEntry, j.cfg.GetMediaExtensions())
	slices.SortFunc(allMediaFiles, util.SortBySizeDesc)

	allSubtitleFiles := j.getAllEntriesOfExtension(rootEntry, j.cfg.GetSubtitleExtensions())
	slices.SortFunc(allSubtitleFiles, util.SortBySizeDesc)

	allSelected := []RenameEntry{}
	allIgnored := []filesystem.DirEntry{}

	for _, mediaFile := range allMediaFiles {
		info := j.GetMediaSeasonAndEpisode(mediaFile.Path)
		if info.Season > 0 && info.Episode > 0 {
			// check if that season and episode is already added
			// since this is sorted by size descending, I am assuming
			// that largest file is correct and already present in the list.
			sameEntries := util.Filter(allSelected, func(re RenameEntry) bool {
				return re.Season == info.Season && re.Episode == info.Episode
			})
			if len(sameEntries) > 0 {
				// add to ignored list
				allIgnored = append(allIgnored, mediaFile)
			} else {
				// add to rename select list
				allSelected = append(allSelected, RenameEntry{
					Media:    &mediaFile,
					Subtitle: nil,
					Season:   info.Season,
					Episode:  info.Episode,
				})
			}
		} else {
			allIgnored = append(allIgnored, mediaFile)
		}
	}

	for _, srtFile := range allSubtitleFiles {
		info := j.GetMediaSeasonAndEpisode(srtFile.Path)
		if info.Season > 0 && info.Episode > 0 {
			// check if a media is there for this file
			targetMediaIndex := -1
			for i, sel := range allSelected {
				if sel.Season == info.Season && sel.Episode == info.Episode {
					targetMediaIndex = i
					break
				}
			}
			if targetMediaIndex > -1 && allSelected[targetMediaIndex].Subtitle == nil {
				allSelected[targetMediaIndex].Subtitle = &srtFile
			} else {
				allIgnored = append(allIgnored, srtFile)
			}
		} else {
			allIgnored = append(allIgnored, srtFile)
		}
	}

	// sort by season and episode
	slices.SortFunc(allSelected, func(a, b RenameEntry) int {
		if a.Season > b.Season {
			return 1
		}
		if a.Season < b.Season {
			return -1
		}
		if a.Episode > b.Episode {
			return 1
		}
		if a.Episode < b.Episode {
			return -1
		}
		return 0
	})

	out := EntriesAndIgnores{
		Selected: allSelected,
		Ignored:  allIgnored,
	}
	return out
}
func (j *JmrRenamer) getAllEntriesOfExtension(rootEntry filesystem.DirEntry, ext []string) []filesystem.DirEntry {

	if rootEntry.IsDirectory {
		allChilds := []filesystem.DirEntry{}
		for _, child := range rootEntry.Children {
			t := j.getAllEntriesOfExtension(child, ext)
			allChilds = append(allChilds, t...)
		}
		return allChilds
	} else {
		e := path.Ext(rootEntry.Path)
		if slices.Contains(ext, e) {
			return []filesystem.DirEntry{
				rootEntry,
			}
		}
		return []filesystem.DirEntry{}
	}
}

func removeSpecialCharacters(inputFilename string) string {
	outputFilename := ""
	for _, ch := range inputFilename {
		re := regexp.MustCompile("[A-Za-z0-9 &]")

		if re.MatchString(string(ch)) {
			outputFilename += string(ch)
		} else {
			outputFilename += " "
		}
	}
	return outputFilename
}

func extractYear(filename *string) int {
	re := regexp.MustCompile(" ?[0-9]{4,} ?")

	if re.MatchString(*filename) {
		matchedYearString := re.FindString(*filename)
		matchedYearIndex := re.FindStringIndex(*filename)
		if len(matchedYearString) == 0 {
			return -1
		}
		yearInt, err := strconv.Atoi(strings.Trim(matchedYearString, " "))
		if err != nil {
			log.Fatal("Cannot convert year string to integer", err)
		}
		*filename = (*filename)[0:matchedYearIndex[0]]

		if yearInt > 9999 || yearInt < 1900 {
			yearInt = -1
		}
		return int(yearInt)
	}
	return -1
}

func removeDoubleWhitespace(str string) string {
	re := regexp.MustCompile(" {2,}")

	singleSpacedStr := re.ReplaceAllString(str, " ")
	return strings.Trim(singleSpacedStr, " ")
}
