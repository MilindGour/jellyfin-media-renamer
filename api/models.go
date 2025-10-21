package api

import (
	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
)

type ConfigResponse struct {
	Version           string                   `json:"version"`
	Port              string                   `json:"port"`
	AllowedExtensions config.AllowedExtensions `json:"allowedExtensions"`
	Source            []DirConfigWithID        `json:"source"`
}

func NewConfigResponse(config *config.Config) *ConfigResponse {
	dirSrcWithID := []DirConfigWithID{}
	for i, cfg := range config.Source {
		dirSrcWithID = append(dirSrcWithID, DirConfigWithID{
			DirConfig: cfg,
			ID:        i,
		})
	}

	return &ConfigResponse{
		Version:           config.Version,
		Port:              config.Port,
		AllowedExtensions: config.AllowedExtensions,
		Source:            dirSrcWithID,
	}
}

type SourcesResponse struct {
	Sources []DirConfigWithID `json:"sources"`
}

func NewSourcesResponse(in []DirConfigWithID) SourcesResponse {
	return SourcesResponse{
		Sources: in,
	}
}

type SourceByIDResponse struct {
	Source  DirConfigWithID  `json:"source"`
	Entries []DirEntryWithID `json:"entries"`
}
type DirEntryWithID struct {
	filesystem.DirEntry
	ID int `json:"id"`
}
type DirConfigWithID struct {
	config.DirConfig
	ID int `json:"id"`
}

func NewSourceByIDResponse(src DirConfigWithID, children []filesystem.DirEntry) SourceByIDResponse {
	out := SourceByIDResponse{
		Source: src,
	}

	for index, child := range children {
		out.Entries = append(out.Entries, DirEntryWithID{
			ID: (index + 1),
			DirEntry: filesystem.DirEntry{
				Name:        child.Name,
				Path:        child.Path,
				Size:        child.Size,
				IsDirectory: child.IsDirectory,
			},
		})
	}

	return out
}

type IdentifyNameRequestItem struct {
	Entry    DirEntryWithID              `json:"entry"`
	Type     mediainfoprovider.MediaType `json:"type"`
	Selected bool                        `json:"selected"`
}
type IdentifyNameRequest []IdentifyNameRequestItem

type IdentifyMediaResponseItem struct {
	SourceDirectory      IdentifyNameRequestItem       `json:"sourceDirectory"`
	IdentifiedMediaName  string                        `json:"identifiedMediaName"`
	IdentifiedMediaYear  int                           `json:"identifiedMediaYear"`
	IdentifiedMediaId    string                        `json:"identifiedMediaId"`
	IdentifiedMediaInfos []mediainfoprovider.MediaInfo `json:"identifiedMediaInfos"`
}
type IdentifyMediaResponse []IdentifyMediaResponseItem

// Next media request and response items are same type
type IdentifyMediaRequestItem IdentifyMediaResponseItem
type IdentifyMediaRequest IdentifyMediaResponse

type RenameMediaRequestItem IdentifyMediaResponseItem
type RenameMediaRequest []RenameMediaRequestItem

type DestinationResponseItem struct {
	Name        string                      `json:"name"`
	Path        string                      `json:"path"`
	Type        mediainfoprovider.MediaType `json:"type"`
	ID          int                         `json:"id"`
	MountPoint  string                      `json:"mount_point"`
	TotalSizeKB int64                       `json:"total_size_kb"`
	FreeSizeKB  int64                       `json:"free_size_kb"`
	UsedSizeKB  int64                       `json:"used_size_kb"`
}

type DestinationResponse []DestinationResponseItem
