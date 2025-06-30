package models

import "fmt"

type MediaType string

const (
	MediaTypeMovie MediaType = "MOVIE"
	MediaTypeTV    MediaType = "TV"
)

// Provides common fields between movie and tv
type MediaInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	YearOfRelease int    `json:"yearOfRelease"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
	MediaId       string `json:"mediaId"`
}

func (m MediaInfo) String() string {
	return fmt.Sprintf("%s(%d) tmdbid-%s", m.Name, m.YearOfRelease, m.MediaId)
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
