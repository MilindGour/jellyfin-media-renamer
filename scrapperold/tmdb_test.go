package scrapperold_test

import (
	"log"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/mock"
	"github.com/MilindGour/jellyfin-media-renamer/models"
	scrapper "github.com/MilindGour/jellyfin-media-renamer/scrapperold"
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

func TestTmdbScrapper_GetProperEpisodeName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		mediaName string
		season    int
		episode   int
		want      string
	}{
		{
			name:      "with single digit season and episode",
			mediaName: "Test",
			season:    1,
			episode:   2,
			want:      "Season 01/Test S01E02",
		},
		{
			name:      "with double digit season and episode",
			mediaName: "Test",
			season:    10,
			episode:   22,
			want:      "Season 10/Test S10E22",
		},
		{
			name:      "with invalid season episode",
			mediaName: "Test",
			season:    -1,
			episode:   -1,
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := scrapper.NewTmdbScrapper()
			got := tm.GetProperEpisodeName(tt.mediaName, tt.season, tt.episode)

			if got != tt.want {
				t.Errorf("GetProperEpisodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbScrapper_GetProperMediaName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		mediaName string
		year      int
		mediaID   string
		want      string
	}{
		{
			name:      "all fields valid",
			mediaName: "Test",
			year:      1992,
			mediaID:   "11",
			want:      "Test (1992) [tmdbid-11]",
		},
		{
			name:      "without year",
			mediaName: "Test",
			year:      0,
			mediaID:   "11",
			want:      "Test [tmdbid-11]",
		},
		{
			name:      "without mediaID",
			mediaName: "Test",
			year:      1992,
			mediaID:   "",
			want:      "Test (1992)",
		},
		{
			name:      "without year and mediaID",
			mediaName: "Test",
			year:      0,
			mediaID:   "",
			want:      "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := scrapper.NewTmdbScrapper()
			got := tm.GetProperMediaName(tt.mediaName, tt.year, tt.mediaID)

			if got != tt.want {
				t.Errorf("GetProperMediaName() = %v, want %v", got, tt.want)
			}
		})
	}
}
