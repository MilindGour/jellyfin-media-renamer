package util_test

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func TestExtractYearFromString(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		in      string
		want    int
		wantErr bool
	}{
		{
			name: "Extract year from a release date string",
			in:   "28 July 2023",
			want: 2023,
		},
		{
			name:    "Throw error",
			in:      "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := util.ExtractYearFromString(tt.in)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExtractYearFromString() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExtractYearFromString() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ExtractYearFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinPaths(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		paths []string
		want  string
	}{
		{
			name:  "Test relative joins",
			paths: []string{"path1", "path2"},
			want:  "path1/path2",
		},
		{
			name:  "Test absolute-relative joins",
			paths: []string{"/path1", "path2"},
			want:  "/path1/path2",
		},
		{
			name:  "Test absolute joins",
			paths: []string{"https://test", "one", "path"},
			want:  "https://test/one/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.JoinPaths(tt.paths...)
			if got != tt.want {
				t.Errorf("JoinPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanFilename(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		inputFilename string
		want          string
		want2         int
	}{
		{
			name:          "Clean simple filename",
			inputFilename: "Simple",
			want:          "Simple",
			want2:         0,
		},
		{
			name:          "Clean file with year",
			inputFilename: "With year 2015",
			want:          "With year",
			want2:         2015,
		},
		{
			name:          "Clean file with some special characters",
			inputFilename: "Special.Characters(2022) [mkv] 1080p",
			want:          "Special Characters",
			want2:         2022,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := util.CleanFilename(tt.inputFilename)
			if got != tt.want {
				t.Errorf("CleanFilename() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("CleanFilename() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestExtractMediaIdFromUrl(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   string
		want string
	}{
		{
			name: "Extract media id from movie url",
			in:   "https://foobar.com/movie/123",
			want: "123",
		},
		{
			name: "Extract media id from complex movie url",
			in:   "https://foobar.com/movie/1234-movie-name",
			want: "1234",
		},
		{
			name: "Extract media id from relative url",
			in:   "/tv/4434-complex-name",
			want: "4434",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ExtractMediaIdFromUrl(tt.in)
			if got != tt.want {
				t.Errorf("ExtractMediaIdFromUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTotalEpisodesFromInfoString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		infoString string
		want       int
	}{
		{
			name:       "Extract total episodes from general string",
			infoString: "1990 • 27 Episodes",
			want:       27,
		},
		{
			name:       "Extract total episodes from erronous string",
			infoString: "13 Episodes",
			want:       13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ExtractTotalEpisodesFromInfoString(tt.infoString)
			if got != tt.want {
				t.Errorf("ExtractTotalEpisodesFromInfoString() = %v, want %v", got, tt.want)
			}
		})
	}
}
