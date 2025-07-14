package scrapper

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
)

func setTestDataForRenameTests() error {
	scrapSearchResponseJSON, err := os.ReadFile("../testdata/tmdb_scrap_search_response.json")
	if err != nil {
		log.Println("Error reading file:", err.Error())
		return err
	}
	state.LastSecondPageAPIResponse = &models.SecondScreenResponse{}
	state.LastConfigSourceByID = &models.ConfigSourceByIDResponse{
		BasePath: "/dummy/test/path",
	}

	return json.Unmarshal(scrapSearchResponseJSON, state.LastSecondPageAPIResponse)
}

func Test_getSingleMovieRenames(t *testing.T) {
	err := setTestDataForRenameTests()
	if err != nil {
		t.Errorf("Cannot initialize test data. %s", err.Error())
	}

	tests := []struct {
		name           string
		id             int
		movieResult    models.MovieResult
		wantNilResult  bool
		wantTotalMedia int
		wantTotalSrt   int
	}{
		{
			name:           "2 media 1 srt test",
			id:             2,
			movieResult:    models.MovieResult{MediaInfo: models.MediaInfo{Name: "Test Movie 1", Description: "Test description 1", YearOfRelease: 1980, MediaID: "1111"}},
			wantNilResult:  false,
			wantTotalMedia: 2,
			wantTotalSrt:   1,
		},
		{
			name:           "1 media 0 srt test",
			id:             7,
			movieResult:    models.MovieResult{MediaInfo: models.MediaInfo{Name: "Test Movie 2", Description: "Test description 2", YearOfRelease: 1999, MediaID: "2222"}},
			wantNilResult:  false,
			wantTotalMedia: 1,
			wantTotalSrt:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSingleMovieRenames(tt.id, tt.movieResult)
			log.Println(got)

			if got == nil && tt.wantNilResult == false {
				t.Errorf("Wanted nil result. But got: %v", got)
			}
			if len(got.MediaRenames) != tt.wantTotalMedia {
				t.Errorf("Wanted %d media results. Got %d", tt.wantTotalMedia, len(got.MediaRenames))
			}
			if len(got.SubtitleRenames) != tt.wantTotalSrt {
				t.Errorf("Wanted %d srt results. Got %d", tt.wantTotalSrt, len(got.SubtitleRenames))
			}
		})
	}
}
