package scrapper

import (
	"encoding/json"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

func setTestDataForRenameTests() error {
	state.LastSecondPageAPIResponse = &models.SecondScreenResponse{}
	state.LastConfigSourceByID = &models.ConfigSourceByIDResponse{
		BasePath: "/dummy/test/path",
	}

	return json.Unmarshal(testdata.ScrapSearchResponseMock, state.LastSecondPageAPIResponse)
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
			name:           "nil movie check",
			id:             1,
			movieResult:    models.MovieResult{MediaInfo: models.MediaInfo{Name: "Test Movie 1", Description: "Test description 1", YearOfRelease: 1980, MediaID: "1111"}},
			wantNilResult:  true,
			wantTotalMedia: 0,
			wantTotalSrt:   0,
		},
		{
			name:           "2 media 1 srt test",
			id:             2,
			movieResult:    models.MovieResult{MediaInfo: models.MediaInfo{Name: "Test Movie 1", Description: "Test description 1", YearOfRelease: 1980, MediaID: "1111"}},
			wantNilResult:  false,
			wantTotalMedia: 1,
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
			s := getMockScrapperClient()
			got := s.getSingleMovieRenames(tt.id, tt.movieResult)

			if got == nil {
				if !tt.wantNilResult {
					t.Errorf("Wanted non-nil result. But got nil.")
				}
			} else {
				// got is not nil
				if tt.wantNilResult {
					t.Errorf("Wanted nil result. But got: %v", got)
				}
				if len(got.RootRenames) != 1 {
					t.Errorf("Wanted only 1 root. Got %d", len(got.RootRenames))
				}
				if len(got.MediaRenames) != tt.wantTotalMedia {
					t.Errorf("Wanted %d media results. Got %d", tt.wantTotalMedia, len(got.MediaRenames))
				}
				if len(got.SubtitleRenames) != tt.wantTotalSrt {
					t.Errorf("Wanted %d srt results. Got %d", tt.wantTotalSrt, len(got.SubtitleRenames))
				}
			}
		})
	}
}

func Test_parseSeasonAndEpisodeNumberFromFilepath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filename string
		want     int
		want2    int
	}{
		{
			filename: "/some/random/test/case/fileS02E04.mp4",
			name:     "Parse SXXEXX",
			want:     2,
			want2:    4,
		},
		{
			filename: "/some/random/test/case/Some Cool TV Show Season 3 - 01.someone.mp4",
			name:     "Parse Season X - XX",
			want:     3,
			want2:    1,
		},
		{
			filename: "/some/random/test/case/Some Cool TV Show Season 3 - 01v2.someone.mp4",
			name:     "Parse Season X - XXvY",
			want:     3,
			want2:    1,
		},
		{
			filename: "/test/cool_show_Season_5_-_02_Other_details.mkv",
			name:     "Parse NNN_Season_X_-_XX",
			want:     5,
			want2:    2,
		},
		{
			filename: "/some/random/test/case/another S1 - 09.mp4",
			name:     "Parse SX - XX",
			want:     1,
			want2:    9,
		},
		{
			filename: "/some/random/test/case/cool_show_Episode_-_06.mp4",
			name:     "Parse XX (No season number) Assume season 1",
			want:     1,
			want2:    6,
		},
		{
			filename: "/some/random/test/case/cool_show_04_Dual Audio_10bit_BD1080p.mp4",
			name:     "Parse XX (No season number) Assume season 1",
			want:     1,
			want2:    4,
		},
		{
			filename: "/some/random/TV Show [1080]/TV_SHOW_-02.mkv",
			name:     "Parse XXXX XX (No season number)",
			want:     1,
			want2:    2,
		},
		{
			filename: "some/random/name/with/no info.mp4",
			name:     "Parse XX (No season number) Assume season 1",
			want:     -1,
			want2:    -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrapperClient()
			got, got2 := s.parseSeasonAndEpisodeNumberFromFilepath(tt.filename)

			if got != tt.want {
				t.Errorf("season parseSeasonAndEpisodeNumberFromFilename() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("episode parseSeasonAndEpisodeNumberFromFilename() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getSingleTVRenames(t *testing.T) {
	err := setTestDataForRenameTests()
	if err != nil {
		t.Errorf("Cannot initialize test data. %s", err.Error())
	}

	tests := []struct {
		name           string
		id             int
		tvResult       models.TVResult
		wantTotalMedia int
		wantTotalSrt   int
		wantMediaName  string
		wantNilResult  bool
	}{
		{
			name:           "Test TV 1",
			id:             9,
			tvResult:       models.TVResult{MediaInfo: models.MediaInfo{Name: "Test A", Description: "Test description", YearOfRelease: 2013, MediaID: "123"}},
			wantTotalMedia: 10,
			wantTotalSrt:   1,
			wantNilResult:  false,
		},
		{
			name:           "Test TV 2",
			id:             21,
			tvResult:       models.TVResult{MediaInfo: models.MediaInfo{Name: "Test B", Description: "Test description", YearOfRelease: 2013, MediaID: "456"}},
			wantTotalMedia: 12,
			wantTotalSrt:   0,
			wantNilResult:  false,
		},
		{
			name:           "Test TV 3",
			id:             34,
			tvResult:       models.TVResult{MediaInfo: models.MediaInfo{Name: "Test C", Description: "Test description", YearOfRelease: 2017, MediaID: "789"}},
			wantTotalMedia: 11,
			wantTotalSrt:   0,
			wantNilResult:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := getMockScrapperClient()
			got := s.getSingleTVRenames(tt.id, tt.tvResult)

			if got == nil {
				if !tt.wantNilResult {
					t.Errorf("getSingleTVRenames() = nil, want non-nil")
				}
			} else {
				if tt.wantNilResult {
					t.Errorf("getSingleTVRenames() = %v, want nil", got)
				}
				if tt.wantTotalMedia != len(got.MediaRenames) {
					t.Errorf("total media getSingleTVRenames() = %d, want %d", len(got.MediaRenames), tt.wantTotalMedia)
				}
				if tt.wantTotalSrt != len(got.SubtitleRenames) {
					t.Errorf("total srt getSingleTVRenames() = %d, want %d", len(got.SubtitleRenames), tt.wantTotalSrt)
				}
			}
		})
	}
}

func getMockScrapperClient() *ScrapperClient {
	s := NewScrapperClient()
	mockConfig := config.NewJmrConfigByData(testdata.ConfigJsonMock)
	s.config = mockConfig

	return s
}
