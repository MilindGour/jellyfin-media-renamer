package api

import (
	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
)

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

type IdentifyNamesRequestItem struct {
	Entry    DirEntryWithID              `json:"entry"`
	Type     mediainfoprovider.MediaType `json:"type"`
	Selected bool                        `json:"selected"`
}
type IdentifyNamesRequest []IdentifyNamesRequestItem

type IdentifyNamesResponseItem struct {
	SourceDirectory      IdentifyNamesRequestItem      `json:"sourceDirectory"`
	IdentifiedMediaName  string                        `json:"identifiedMediaName"`
	IdentifiedMediaYear  int                           `json:"identifiedMediaYear"`
	IdentifiedMediaInfos []mediainfoprovider.MediaInfo `json:"identifiedMediaInfos"`
}
type IdentifyNamesResponse []IdentifyNamesResponseItem
