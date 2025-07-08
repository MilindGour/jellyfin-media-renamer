// Package scrapper contains functions that deals with web scrapping.
package scrapper

import (
	"github.com/MilindGour/jellyfin-media-renamer/models"
)

type Scrapper interface {
	GetSearchableString(models.ClearFileEntry) string
	SearchMovie(models.ClearFileEntry) ([]models.MovieResult, error)
	SearchTV(models.ClearFileEntry) ([]models.TVResult, error)
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
	// TODO: implement this function
	panic("getMovieRenames is not implemented")
}

func getSingleTVRenames(id int, tvResult models.TVResult) *models.TVRenameResult {
	// TODO: implement this function
	panic("getTVRenames is not implemented")
}
