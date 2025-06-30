package tests

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
)

func TestScrapSeasonInfoList(t *testing.T) {
	client := scrapper.NewTmdbScrapper()
	result := client.ScrapSeasonInfoList("4327")

	if len(result) == 0 {
		t.Error("Expected output to have 2 length. Found 0.")
	}
}
