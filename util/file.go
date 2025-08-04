package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/models"
)

var nextFileID int = 0

func getNextFileID() int {
	nextFileID += 1
	return nextFileID
}

func ResetNextFileID() {
	nextFileID = 0
}

func GetDirectoryEntries(path string, allowedExtensions []string) ([]models.DirectoryEntry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Cannot read directory. ", err))
	}

	out := []models.DirectoryEntry{}

	for _, entry := range dirEntries {
		curEntry := models.DirectoryEntry{}
		eInfo, err2 := entry.Info()

		if err2 != nil {
			return nil, errors.New(fmt.Sprint("Error reading entry info. ", err))
		}

		curEntry.ID = getNextFileID()
		curEntry.Name = eInfo.Name()
		curEntry.Path = fmt.Sprintf("%s/%s", path, curEntry.Name)
		curEntry.Size = eInfo.Size()
		curEntry.IsDirectory = eInfo.IsDir()

		if curEntry.IsDirectory {
			childEntries, err := GetDirectoryEntries(JoinPaths(path, curEntry.Name), allowedExtensions)
			if err != nil {
				return nil, err
			}
			// calculate size of directory as OS does not provide correct size
			dirsize := int64(0)
			for _, sub := range childEntries {
				dirsize += sub.Size
			}
			curEntry.Size = dirsize
			curEntry.Children = childEntries
		} else {
			curEntry.Children = nil

			// if the file is neither media nor subtitle, skip it
			fileExtension := filepath.Ext(curEntry.Path)
			if !HasItem(allowedExtensions, func(x string) bool {
				return fileExtension == x
			}) {
				continue
			}
		}

		out = append(out, curEntry)
	}

	return out, nil
}

func GetConfigFilename() string {
	if IsProduction() {
		configFilename := fmt.Sprintf("%s/.config/jmr/config.json", os.Getenv("HOME"))
		return configFilename
	}
	return "config.json"
}

// FilterSubtitleFileEntries recursively filters the subtitle file child entries from a given directory.
func FilterSubtitleFileEntries(in models.DirectoryEntry, subtitleExtensions []string) []models.DirectoryEntry {
	return FilterDirectoryEntries(in, fileExtensionFilterFunction(subtitleExtensions))
}

// FilterVideoFileEntries recursively filters the video file child entries from a given directory.
func FilterVideoFileEntries(in models.DirectoryEntry, videoExtensions []string) []models.DirectoryEntry {
	return FilterDirectoryEntries(in, fileExtensionFilterFunction(videoExtensions))
}

func fileExtensionFilterFunction(extensions []string) func(models.DirectoryEntry) bool {
	return func(de models.DirectoryEntry) bool {
		for _, ext := range extensions {
			if strings.HasSuffix(de.Path, ext) {
				return true
			}
		}
		return false
	}
}

func FilterDirectoryEntries(in models.DirectoryEntry, predicate func(models.DirectoryEntry) bool) []models.DirectoryEntry {
	out := []models.DirectoryEntry{}
	if !in.IsDirectory {
		result := predicate(in)
		if result {
			t := models.DirectoryEntry{
				ID:   in.ID,
				Name: in.Name,
				Path: in.Path,
				Size: in.Size,
			}
			out = append(out, t)
		}
	} else {
		for _, childDir := range in.Children {
			childResult := FilterDirectoryEntries(childDir, predicate)
			out = append(out, childResult...)
		}
	}
	return out
}

func SortByFileSizeDescending(a, b models.DirectoryEntry) int {
	return int(b.Size - a.Size)
}

func SortBySeasonAndEpisodeNumbers(a, b models.MediaPathRename) int {
	re := regexp.MustCompile(`S(\d{2})E(\d{2})`)
	matchA := re.FindStringSubmatch(a.NewPath)
	matchB := re.FindStringSubmatch(b.NewPath)

	seasonA, err1 := strconv.Atoi(matchA[1])
	episodeA, err2 := strconv.Atoi(matchA[2])
	if err1 != nil || err2 != nil {
		return -1
	}

	seasonB, err1 := strconv.Atoi(matchB[1])
	episodeB, err2 := strconv.Atoi(matchB[2])
	if err1 != nil || err2 != nil {
		return 1
	}

	if seasonA > seasonB {
		return 1
	} else if seasonA < seasonB {
		return -1
	}
	if episodeA > episodeB {
		return 1
	}
	if episodeB > episodeA {
		return -1
	}
	return 0
}
