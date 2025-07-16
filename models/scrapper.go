// Package models contains all the data structures used in the app.
package models

import "fmt"

type MediaType string

const (
	MediaTypeMovie MediaType = "MOVIE"
	MediaTypeTV    MediaType = "TV"
)

// MediaInfo provides common fields between movie and tv
type MediaInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	YearOfRelease int    `json:"yearOfRelease"`
	ThumbnailURL  string `json:"thumbnailUrl"`
	MediaID       string `json:"mediaId"`
}

func (m MediaInfo) String() string {
	return fmt.Sprintf("%s(%d) tmdbid-%s", m.Name, m.YearOfRelease, m.MediaID)
}

type MovieResult struct {
	MediaInfo
}

type TVResult struct {
	MediaInfo
	TotalSeasons int          `json:"totalSeasons"`
	Seasons      []SeasonInfo `json:"seasons"`
}

type SeasonInfo struct {
	Number        int `json:"number"`
	TotalEpisodes int `json:"totalEpisodes"`
}

type ScrapSearchRequest struct {
	CleanFilenameEntries map[int]ClearFileEntry `json:"cleanFilenameEntries"`
	MediaTypes           map[int]MediaType      `json:"mediaTypes"`
}

type ScrapSearchConfirmRequest struct {
	MoviesInfo map[int]MovieResult `json:"moviesInfo"`
	TVsInfo    map[int]TVResult    `json:"tvsInfo"`
}

type ScrapSearchRenameResult struct {
	MovieRenameResults []MovieRenameResult `json:"movieRenameResults"`
	TVRenameResults    []TVRenameResult    `json:"tvRenameResults"`
}

type MovieRenameResult struct {
	MovieResult       MovieResult       `json:"movieResult"`
	RootRenames       []MediaPathRename `json:"rootRenames"`
	MediaRenames      []MediaPathRename `json:"mediaRenames"`
	SubtitleRenames   []MediaPathRename `json:"subtitleRenames"`
	IgnoredMediaPaths []string          `json:"ignoredMediaPaths"`
}

type TVRenameResult struct {
	TVResult          TVResult          `json:"tvResult"`
	RootRenames       []MediaPathRename `json:"rootRenames"`
	MediaRenames      []MediaPathRename `json:"mediaRenames"`
	SubtitleRenames   []MediaPathRename `json:"subtitleRenames"`
	IgnoredMediaPaths []string          `json:"ignoredMediaPaths"`
}

func (r TVRenameResult) String() string {
	out := ""
	out += fmt.Sprintf("TV: %s\nMedia Renames:\n", r.TVResult.Name)
	for _, mediaRename := range r.MediaRenames {
		out += fmt.Sprintf("\nOLD: %s\nNEW: %s\n", mediaRename.OldPath, mediaRename.NewPath)
	}
	if len(r.SubtitleRenames) > 0 {
		out += "\n\nSubtitles:\n"
		for _, srt := range r.SubtitleRenames {
			out += fmt.Sprintf("\nOLD: %s\nNEW: %s\n", srt.OldPath, srt.NewPath)
		}
	}
	if len(r.IgnoredMediaPaths) > 0 {
		out += "\n\nIgnored files:\n"
		for _, ignoredFilepath := range r.IgnoredMediaPaths {
			out += fmt.Sprintf("\n-> %s", ignoredFilepath)
		}
		out += "\n"
	}
	return out
}

type MediaPathRename struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
}
