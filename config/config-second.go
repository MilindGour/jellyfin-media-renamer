package config

import (
	"errors"
	"fmt"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func PopulateSecondScreenResponse(selectedIds []int) (*models.SecondScreenResponse, error) {
	response := models.SecondScreenResponse{
		MovieResults: map[int][]models.MovieResult{},
		TVResults:    map[int]models.TVResult{},
	}

	// select all the directories using selectedIds
	for _, selectedId := range selectedIds {
		entry := util.Filter(state.LastConfigSourceById, func(x models.DirectoryEntry) bool {
			return x.Id == int(selectedId)
		})
		if len(entry) > 0 {
			response.SelectedDirEntries = append(response.SelectedDirEntries, entry[0])
		} else {
			errmsg := fmt.Sprintf("Cannot find entry with id: %d. Probably you need to call the config/source/:id api first.", selectedId)
			return nil, errors.New(errmsg)
		}
	}
	response.Success = len(response.SelectedDirEntries) > 0
	populateCleanFilenames(&response)

	return &response, nil
}

func populateCleanFilenames(in *models.SecondScreenResponse) {
	in.CleanFilenameEntries = map[int]models.ClearFileEntry{}
	for _, dirEntry := range in.SelectedDirEntries {
		cleanName, year := util.CleanFilename(dirEntry.Name)
		in.CleanFilenameEntries[dirEntry.Id] = models.ClearFileEntry{
			Name: cleanName,
			Year: year,
		}
	}
}
