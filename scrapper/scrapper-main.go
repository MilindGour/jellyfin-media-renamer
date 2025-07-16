// Package scrapper contains functions that deals with web scrapping.
package scrapper

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type Scrapper interface {
	GetSearchableString(models.ClearFileEntry) string
	SearchMovie(models.ClearFileEntry) ([]models.MovieResult, error)
	SearchTV(models.ClearFileEntry) ([]models.TVResult, error)
	GetProperMediaName(mediaName string, year int, mediaID string) string
	GetProperEpisodeName(mediaName string, season, episode int) string
}

func NewPathRename(oldPath, newPath string) *models.MediaPathRename {
	return &models.MediaPathRename{
		OldPath: oldPath,
		NewPath: newPath,
	}
}

func ValidateScrapConfirmRequest(in models.ScrapSearchConfirmRequest) bool {
	return len(in.MoviesInfo) > 0 || len(in.TVsInfo) > 0
}

func GetMediaRenames(in models.ScrapSearchConfirmRequest, selectedDirectoryEntries []models.DirectoryEntry) *models.ScrapSearchRenameResult {
	out := &models.ScrapSearchRenameResult{}
	// process movies
	for movieDirID, movieResult := range in.MoviesInfo {
		movieRenameResult := getSingleMovieRenames(movieDirID, movieResult)
		out.MovieRenameResults = append(out.MovieRenameResults, *movieRenameResult)
	}

	// process tvs
	for tvDirID, tvResult := range in.TVsInfo {
		tvRenameResult := getSingleTVRenames(tvDirID, tvResult)
		out.TVRenameResults = append(out.TVRenameResults, *tvRenameResult)
	}

	return out
}

func getSingleMovieRenames(id int, movieResult models.MovieResult) *models.MovieRenameResult {
	oldBase := state.LastConfigSourceByID.BasePath
	newBase := fmt.Sprintf("%s/%s", oldBase, "jmr_renames")

	// TODO: Once new scrap client is implemented, add a logic to select appropriate client.
	var scrapperClient Scrapper = NewTmdbScrapper()
	renameString := scrapperClient.GetProperMediaName(movieResult.Name, movieResult.YearOfRelease, movieResult.MediaID)

	out := &models.MovieRenameResult{
		MovieResult: movieResult,
	}

	// Get the selectedDirectoryEntry by id
	entries := util.Filter(state.LastSecondPageAPIResponse.SelectedDirEntries, func(de models.DirectoryEntry) bool {
		return de.ID == id
	})
	if len(entries) == 0 {
		return nil
	}
	targetDirEntry := entries[0]

	// Change the name of parentDirectory
	movieRootOld := fmt.Sprintf("%s/%s", oldBase, targetDirEntry.Name)
	movieRootNew := fmt.Sprintf("%s/%s", newBase, renameString)
	out.RootRenames = append(out.RootRenames, *NewPathRename(movieRootOld, movieRootNew))

	// Get main video file
	videoEntries := util.FilterVideoFileEntries(targetDirEntry)
	if len(videoEntries) > 0 {
		slices.SortFunc(videoEntries, util.SortByFileSizeDescending)

		// Take the largest size video as the main video file.
		movieFileOld := videoEntries[0].Path
		movieFileExt := path.Ext(videoEntries[0].Name)
		movieFilenameNew := fmt.Sprintf("%s%s", renameString, movieFileExt)
		movieFileNew := util.JoinPaths(movieRootNew, movieFilenameNew)
		out.MediaRenames = append(out.MediaRenames, *NewPathRename(movieFileOld, movieFileNew))

		// Put rest of the video file paths in IgnoredMediapaths so that frontend can show.
		if len(videoEntries) > 1 {
			for i := 1; i < len(videoEntries); i++ {
				out.IgnoredMediaPaths = append(out.IgnoredMediaPaths, videoEntries[i].Path)
			}
		}
	}

	// Get subtitle
	subtitleEntries := util.FilterSubtitleFileEntries(targetDirEntry)
	if len(subtitleEntries) > 0 {
		slices.SortFunc(subtitleEntries, util.SortByFileSizeDescending)
		subtitleFileOld := subtitleEntries[0].Path
		subtitleFileExt := path.Ext(subtitleEntries[0].Name)
		subtitleFilenameNew := fmt.Sprintf("%s%s", renameString, subtitleFileExt)
		subtitleFileNew := util.JoinPaths(movieRootNew, subtitleFilenameNew)
		out.SubtitleRenames = append(out.SubtitleRenames, *NewPathRename(subtitleFileOld, subtitleFileNew))
	}

	return out
}

func getSingleTVRenames(id int, tvResult models.TVResult) *models.TVRenameResult {
	oldBase := state.LastConfigSourceByID.BasePath
	newBase := fmt.Sprintf("%s/%s", oldBase, "jmr_renames")

	// TODO: Once new scrap client is implemented, add a logic to select appropriate client.
	var scrapperClient Scrapper = NewTmdbScrapper()

	// properMediaName is the root directory name of the tv show
	properMediaName := scrapperClient.GetProperMediaName(tvResult.Name, tvResult.YearOfRelease, tvResult.MediaID)

	out := &models.TVRenameResult{
		TVResult: tvResult,
	}

	// Get the selectedDirectoryEntry by id
	entries := util.Filter(state.LastSecondPageAPIResponse.SelectedDirEntries, func(de models.DirectoryEntry) bool {
		return de.ID == id
	})
	if len(entries) == 0 {
		return nil
	}
	targetDirEntry := entries[0]

	// Change the name of parentDirectory
	tvRootOld := fmt.Sprintf("%s/%s", oldBase, targetDirEntry.Name)
	tvRootNew := fmt.Sprintf("%s/%s", newBase, properMediaName)
	out.RootRenames = append(out.RootRenames, *NewPathRename(tvRootOld, tvRootNew))

	// Get video files
	videoEntries := util.FilterVideoFileEntries(targetDirEntry)
	seasonMap := map[int][]int{}
	if len(videoEntries) > 0 {
		slices.SortFunc(videoEntries, util.SortByFileSizeDescending)
		for _, ve := range videoEntries {
			tvFileOld := ve.Path
			tvFileExt := path.Ext(ve.Name)
			season, episode := parseSeasonAndEpisodeNumberFromFilepath(ve.Path)
			if (season == -1 && episode == -1) || isSeasonEpisodeAlreadyPresent(seasonMap, season, episode) {
				out.IgnoredMediaPaths = append(out.IgnoredMediaPaths, ve.Path)
				continue
			}

			tvFilenameNew := fmt.Sprintf("%s%s", scrapperClient.GetProperEpisodeName(tvResult.Name, season, episode), tvFileExt)
			tvFileNew := util.JoinPaths(tvRootNew, tvFilenameNew)
			out.MediaRenames = append(out.MediaRenames, *NewPathRename(tvFileOld, tvFileNew))
		}

		// sort the media files by season and episodes.
		slices.SortFunc(out.MediaRenames, util.SortBySeasonAndEpisodeNumbers)
	}

	// Get srt files
	srtEntries := util.FilterSubtitleFileEntries(targetDirEntry)
	if len(srtEntries) > 0 {
		slices.SortFunc(srtEntries, util.SortByFileSizeDescending)
		for _, ve := range srtEntries {
			srtFileOld := ve.Path
			srtFileExt := path.Ext(ve.Name)
			season, episode := parseSeasonAndEpisodeNumberFromFilepath(ve.Path)
			if season == -1 && episode == -1 {
				continue
			}

			srtFilenameNew := fmt.Sprintf("%s%s", scrapperClient.GetProperEpisodeName(tvResult.Name, season, episode), srtFileExt)
			srtFileNew := util.JoinPaths(tvRootNew, srtFilenameNew)
			out.SubtitleRenames = append(out.SubtitleRenames, *NewPathRename(srtFileOld, srtFileNew))
		}

		// sort the srt files by season and episodes.
		slices.SortFunc(out.SubtitleRenames, util.SortBySeasonAndEpisodeNumbers)
	}

	return out
}

// parseSeasonAndEpisodeNumberFromFilepath returns season and episode numbers with the given filename.
// It returns -1, -1 if cannot find any pattern.
func parseSeasonAndEpisodeNumberFromFilepath(filepath string) (int, int) {
	filepathWithoutExtension, _ := strings.CutSuffix(filepath, path.Ext(filepath))
	in := strings.ToLower(filepathWithoutExtension)
	in += "_"

	testREs := []*regexp.Regexp{
		regexp.MustCompile(`s(\d{2})e(\d{2})`),                     // SXXEXX
		regexp.MustCompile(`season[ \-_]+(\d{1,2})[ \-_]+(\d{2})`), // Season X - XX
		regexp.MustCompile(`s(\d+)[ \-_]+(\d+)`),                   // SX - XX
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
					return -1, -1
				}
				return 1, episode
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
				return season, episode
			}
		}
	}

	return -1, -1
}

func isSeasonEpisodeAlreadyPresent(seasonMap map[int][]int, season, episode int) bool {
	_, seasonFound := seasonMap[season]
	if !seasonFound {
		seasonMap[season] = []int{episode}
		return false
	}

	if slices.Contains(seasonMap[season], episode) {
		// already present
		return true
	}

	seasonMap[season] = append(seasonMap[season], episode)
	return false
}
