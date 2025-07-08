// Package scrapper contains functions that deals with web scrapping.
package scrapper

import (
	"fmt"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type Scrapper interface {
	GetSearchableString(models.ClearFileEntry) string
	SearchMovie(models.ClearFileEntry) ([]models.MovieResult, error)
	SearchTV(models.ClearFileEntry) ([]models.TVResult, error)
	GetRenameString(mediaName string, year int, mediaID string) string
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
	newBase := fmt.Sprintf("%s/%s", oldBase, "__jmr_temp_renames")

	// TODO: Once new scrap client is implemented, add a logic to select appropriate client.
	var scrapperClient Scrapper = NewTmdbScrapper()
	renameString := scrapperClient.GetRenameString(movieResult.Name, movieResult.YearOfRelease, movieResult.MediaID)

	out := &models.MovieRenameResult{
		MovieResult:      movieResult,
		MediaPathRenames: []models.MediaPathRename{},
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
	out.MediaPathRenames = append(out.MediaPathRenames, *NewPathRename(movieRootOld, movieRootNew))

	return out
}

func getSingleTVRenames(id int, tvResult models.TVResult) *models.TVRenameResult {
	// TODO: implement this function
	panic("getTVRenames is not implemented")
}
