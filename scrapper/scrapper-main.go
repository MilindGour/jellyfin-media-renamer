package scrapper

import (
	"github.com/MilindGour/jellyfin-media-renamer/models"
)

type Scrapper interface {
	GetSearchableString(models.ClearFileEntry) string
	SearchMovie(models.ClearFileEntry) []models.MovieResult
	SearchTV(models.ClearFileEntry) []models.TVResult
}
