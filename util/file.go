package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/models"
)

func GetConfigFileContents() ([]byte, error) {
	data, err := os.ReadFile("config.json")
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

		curEntry.Name = eInfo.Name()
		curEntry.Size = eInfo.Size()
		curEntry.IsDirectory = eInfo.IsDir()
		curEntry.Children = nil

		out = append(out, curEntry)
	}

	return out, nil
}
