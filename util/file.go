package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/models"
)

var nextFileId uint = 0

func getNextFileId() uint {
	nextFileId += 1
	return nextFileId
}

func GetConfigFileContents() ([]byte, error) {
	data, err := os.ReadFile(GetConfigFilename())
	return data, err
}

func GetDirectoryEntries(path string) ([]models.DirectoryEntry, error) {
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

		curEntry.Id = getNextFileId()
		curEntry.Name = eInfo.Name()
		curEntry.Size = eInfo.Size()
		curEntry.IsDirectory = eInfo.IsDir()

		if curEntry.IsDirectory {
			childEntries, err := GetDirectoryEntries(JoinPaths(path, curEntry.Name))
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
