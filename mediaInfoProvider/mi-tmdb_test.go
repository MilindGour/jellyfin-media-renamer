package mediainfoprovider

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
)

func TestTmdbMIProvider_SearchMovies(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		term string
		want []MovieResult
	}{
		{
			name: "Movie search with 2 results",
			term: "test movie",
			want: []MovieResult{
				{MediaInfo: MediaInfo{Name: "Blade Runner 2049", Description: "Test content of movie 1", YearOfRelease: 2017, ThumbnailURL: "https://www.test.com/pic/1.jpg", MediaID: "335984"}},
				{MediaInfo: MediaInfo{Name: "Test Movie 2", Description: "Test content for movie 2", YearOfRelease: 2022, ThumbnailURL: "https://www.test.com/pic/2.jpg", MediaID: "954126"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewMockTmdbMIProvider()
			got := tm.SearchMovies(tt.term)

			if len(got) != len(tt.want) {
				t.Errorf("SearchMovies() = len %d, want %d", len(got), len(tt.want))
				return
			}
			for i, gotMovie := range got {
				wantMovieStr := tt.want[i].String()
				gotMovieStr := gotMovie.String()
				if gotMovieStr != wantMovieStr {
					t.Errorf("SearchMovies()=\ngot %s,\nwant %s", gotMovieStr, wantMovieStr)
				}
			}
		})
	}
}

func TestTmdbMIProvider_parseScrapResultListToMediaInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   scrapper.ScrapResultList
		want []MediaInfo
	}{
		{
			name: "Parse map containing all variants",
			in: scrapper.ScrapResultList{
				{"name": "\n\tMovie 1\n\n", "description": "Description 1", "yearOfRelease": "20 December 2011", "thumbnailUrl": "https://mock.url/1.png", "mediaId": "movie/335984-blade-runner-2049"},
				{"name": "\n\t    Movie 2\n\n", "description": "\n\n\t\tComplex description", "yearOfRelease": "", "thumbnailUrl": "https://mock.url/2.png", "mediaId": "https://mock-url.net/tv/12345"},
			},
			want: []MediaInfo{
				MediaInfo{Name: "Movie 1", Description: "Description 1", YearOfRelease: 2011, ThumbnailURL: "https://mock.url/1.png", MediaID: "335984"},
				MediaInfo{Name: "Movie 2", Description: "Complex description", YearOfRelease: 0, ThumbnailURL: "https://mock.url/2.png", MediaID: "12345"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTmdbMIProvider()
			got := tm.parseScrapResultListToMediaInfo(tt.in)

			if len(got) != len(tt.want) {
				t.Errorf("parseScrapResultListToMovieResult() = length %d, wanted %d", len(got), len(tt.want))
				return
			}
			for i, gotMovie := range got {
				gotStr := gotMovie.String()
				wantStr := tt.want[i].String()

				if gotStr != wantStr {
					t.Errorf("parseScrapResultListToMovieResult() = \nGot= %v,\nWanted= %v", gotStr, wantStr)
				}
			}
		})
	}
}

func TestTmdbMIProvider_extraMediaId(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   string
		want string
	}{
		{
			name: "Normal media id extraction",
			in:   "/movie/12345-movie-name",
			want: "12345",
		},
		{
			name: "Complex media id extraction",
			in:   "https://absolute.url/movie/43290-movie-name-2024",
			want: "43290",
		},
		{
			name: "Empty string",
			in:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTmdbMIProvider()
			got := tm.extraMediaId(tt.in)

			if got != tt.want {
				t.Errorf("extraMediaId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_extractYear(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   string
		want int
	}{
		{
			name: "Normal extraction",
			in:   "31 February 2013",
			want: 2013,
		},
		{
			name: "Does not contain year",
			in:   "7 February",
			want: 0,
		},
		{
			name: "Empty string test",
			in:   "",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTmdbMIProvider()
			got := tm.extractYear(tt.in)
			if got != tt.want {
				t.Errorf("extractYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_trimString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		in   string
		want string
	}{
		{
			name: "Space trim",
			in:   "   Test Str  ",
			want: "Test Str",
		},
		{
			name: "Tabbed trim",
			in:   "\t\tNoice\t\t\t",
			want: "Noice",
		},
		{
			name: "Mixed complex trim",
			in:   "\r\n\t\t  This one\n\n\t\t",
			want: "This one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTmdbMIProvider()
			got := tm.trimString(tt.in)

			if got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}
