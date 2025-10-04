package renamer

import (
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
)

type Renamer interface {
	GetMediaNameAndYear(rawFilename string) MediaNameAndYear
	GetMediaSeasonAndEpisode(filePath string) MediaSeasonAndEpisode
	SelectEntriesForRename(rootEntry filesystem.DirEntry, mediaType mediainfoprovider.MediaType) EntriesAndIgnores
	ConfirmEntriesForRename(entries RenameMediaConfirmRequest) (int, error)
}

type MediaNameAndYear struct {
	Name string
	Year int
}

type MediaSeasonAndEpisode struct {
	Season  int
	Episode int
}

type RenameEntry struct {
	Media    *filesystem.DirEntry `json:"media"`
	Subtitle *filesystem.DirEntry `json:"subtitle"`
	Season   int                  `json:"season,omitempty"`
	Episode  int                  `json:"episode,omitempty"`
}

type EntriesAndIgnores struct {
	Selected []RenameEntry         `json:"selected"`
	Ignored  []filesystem.DirEntry `json:"ignored"`
}

type RenameMediaResponseItem struct {
	Info  mediainfoprovider.MediaInfo `json:"info"`
	Type  mediainfoprovider.MediaType `json:"type"`
	Entry filesystem.DirEntry         `json:"entry"`
	EntriesAndIgnores
}

type RenameMediaResponse []RenameMediaResponseItem

type RenameMediaConfirmRequest []RenameMediaResponseItem
