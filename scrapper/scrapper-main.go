// Package scrapper contains functions that deals with web scrapping.
package scrapper

import (
	"fmt"
	"slices"

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
	newBase := fmt.Sprintf("%s/%s", oldBase, "jmr_renames")

	// TODO: Once new scrap client is implemented, add a logic to select appropriate client.
	var scrapperClient Scrapper = NewTmdbScrapper()
	renameString := scrapperClient.GetRenameString(movieResult.Name, movieResult.YearOfRelease, movieResult.MediaID)

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
		for _, ve := range videoEntries {
			movieFileOld := ve.Path
			movieFileExt := util.GetFileExtension(ve.Name)
			movieFilenameNew := fmt.Sprintf("%s%s", renameString, movieFileExt)
			movieFileNew := util.JoinPaths(movieRootNew, movieFilenameNew)
			out.MediaRenames = append(out.MediaRenames, *NewPathRename(movieFileOld, movieFileNew))
		}
	}

	// Get subtitle
	subtitleEntries := util.FilterSubtitleFileEntries(targetDirEntry)
	if len(subtitleEntries) > 0 {
		slices.SortFunc(subtitleEntries, util.SortByFileSizeDescending)
		subtitleFileOld := subtitleEntries[0].Path
		subtitleFileExt := util.GetFileExtension(subtitleEntries[0].Name)
		subtitleFilenameNew := fmt.Sprintf("%s%s", renameString, subtitleFileExt)
		subtitleFileNew := util.JoinPaths(movieRootNew, subtitleFilenameNew)
		out.SubtitleRenames = append(out.SubtitleRenames, *NewPathRename(subtitleFileOld, subtitleFileNew))
	}

	return out
}

func getSingleTVRenames(id int, tvResult models.TVResult) *models.TVRenameResult {
	// TODO: implement this function
	panic("getTVRenames is not implemented")
}
