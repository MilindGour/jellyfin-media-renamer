package mediainfoprovider

import "fmt"

type MediaType string

const (
	MediaTypeMovie MediaType = "MOVIE"
	MediaTypeTV    MediaType = "TV"
)

type MediaInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	YearOfRelease int    `json:"yearOfRelease"`
	ThumbnailURL  string `json:"thumbnailUrl"`
	MediaID       string `json:"mediaId"`
}

func (m MediaInfo) String() string {
	return fmt.Sprintf("[MediaInfo Name=%s, Description=%s, YearOfRelease=%d, ThumbnailURL=%s, MediaID=%s]", m.Name, m.Description, m.YearOfRelease, m.ThumbnailURL, m.MediaID)
}

type MovieResult struct {
	MediaInfo
}

func (m MovieResult) String() string {
	return fmt.Sprintf("[MovieResult %v]", m.MediaInfo)
}

type TVResult struct {
	MediaInfo
	TotalSeasons int          `json:"totalSeasons"`
	Seasons      []SeasonInfo `json:"seasons"`
}

func (t *TVResult) String() string {
	return fmt.Sprintf("[TVResult %s, SS=%d, %v]", t.MediaInfo, t.TotalSeasons, t.Seasons)
}

type SeasonInfo struct {
	Number        int `json:"number"`
	TotalEpisodes int `json:"totalEpisodes"`
}

func (s *SeasonInfo) String() string {
	return fmt.Sprintf("[SeasonInfo S=%2d, E=%2d]", s.Number, s.TotalEpisodes)
}
