package scrapper_test

import (
	"log"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/mock"
	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
)

func TestTmdbScrapper_GetSearchableString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   models.ClearFileEntry
		want string
	}{
		{
			name: "Get search string with year",
			in:   models.ClearFileEntry{Name: "Test", Year: 2012},
			want: "Test y:2012",
		},
		{
			name: "Get search string without year",
			in:   models.ClearFileEntry{Name: "Test", Year: 0},
			want: "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := scrapper.NewTmdbScrapper()
			got := tm.GetSearchableString(tt.in)
			if got != tt.want {
				t.Errorf("GetSearchableString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockServer(t *testing.T) {
	tmdbClient, mockServer := mock.NewMockTmdbClient()
	defer mockServer.Close()

	log.Println("Mock TMDB server url:", tmdbClient.BaseURL)
	result, err := tmdbClient.SearchMovie(models.ClearFileEntry{
		Name: "Airplane",
		Year: 1980,
	})
	if err != nil {
		t.Error("Cannot call SearchMovie. Err =", err.Error())
	}

	log.Println("Total results for Airplane:", len(result))
}
