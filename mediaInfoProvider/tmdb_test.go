package mediainfoprovider

import (
	"reflect"
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
			got := tm.SearchMovies(tt.term, 0)

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
				{Name: "Movie 1", Description: "Description 1", YearOfRelease: 2011, ThumbnailURL: "https://mock.url/1.png", MediaID: "335984"},
				{Name: "Movie 2", Description: "Complex description", YearOfRelease: 0, ThumbnailURL: "https://mock.url/2.png", MediaID: "12345"},
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

func TestTmdbMIProvider_extractMediaId(t *testing.T) {
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
			got := tm.extractMediaId(tt.in)

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

func TestNewTmdbMIProvider(t *testing.T) {
	tests := []struct {
		name string
		want *TmdbMIProvider
	}{
		{
			name: "Normal instance create",
			want: &TmdbMIProvider{
				baseUrl:  "https://www.themoviedb.org",
				scrapper: scrapper.NewGoQueryScrapper(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTmdbMIProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTmdbMIProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMockTmdbMIProvider(t *testing.T) {
	tests := []struct {
		name string
		want *TmdbMIProvider
	}{
		{
			name: "Test mock instance create",
			want: &TmdbMIProvider{
				baseUrl:  "tmdb",
				scrapper: scrapper.NewMockGoQueryScrapper(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMockTmdbMIProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMockTmdbMIProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_getParsedMediaInfoListFromUrl(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		url          string
		itemSelector string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []MediaInfo
	}{
		{
			name: "Parse mock url to mediaInfo",
			fields: fields{
				baseUrl:  "tmdb",
				scrapper: scrapper.NewMockGoQueryScrapper(),
			},
			args: args{
				url:          "tmdb/search/movie?query=test+movie",
				itemSelector: ".search_results.movie .card",
			},
			want: []MediaInfo{
				{Name: "Blade Runner 2049", Description: "Test content of movie 1", YearOfRelease: 2017, ThumbnailURL: "https://www.test.com/pic/1.jpg", MediaID: "335984"},
				{Name: "Test Movie 2", Description: "Test content for movie 2", YearOfRelease: 2022, ThumbnailURL: "https://www.test.com/pic/2.jpg", MediaID: "954126"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.getParsedMediaInfoListFromUrl(tt.args.url, tt.args.itemSelector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmdbMIProvider.getParsedMediaInfoListFromUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_SearchTVShows(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		term string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []TVResult
	}{
		{
			name: "Parse TV Show",
			fields: fields{
				baseUrl:  "tmdb",
				scrapper: scrapper.NewMockGoQueryScrapper(),
			},
			args: args{
				term: "Test TV",
			},
			want: []TVResult{
				{
					MediaInfo:    MediaInfo{Name: "Test TV 1", Description: "Test show 1 description", YearOfRelease: 2012, ThumbnailURL: "test_poster_1", MediaID: "1396"},
					TotalSeasons: 2,
					Seasons: []SeasonInfo{
						{1, 9},
						{2, 18},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.SearchTVShows(tt.args.term, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmdbMIProvider.SearchTVShows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_getSearchItemFieldmap(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "Should return tmdb specific fields",
			fields: fields{
				baseUrl:  "tmdb",
				scrapper: scrapper.NewMockGoQueryScrapper(),
			},
			want: map[string]string{
				"name":          ".title h2",
				"subname":       ".title h2 span.title",
				"description":   ".overview",
				"yearOfRelease": ".title .release_date",
				"thumbnailUrl":  ".image img.poster[src]",
				"mediaId":       ".title a.result[href]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.getSearchItemFieldmap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmdbMIProvider.getSearchItemFieldmap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_getSeasonInformation(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		mediaID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []SeasonInfo
	}{
		{
			"Get season information from mock",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{"1396"},
			[]SeasonInfo{{1, 9}, {2, 18}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			got := tr.getSeasonInformation(tt.args.mediaID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmdbMIProvider.getSeasonInformation() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_extractSeasonNumber(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		in string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"From custom url",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{"/tv/1234-test-tv-1/season/56"},
			56,
		},
		{
			"From short url",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{"/tv/8900/season/10"},
			10,
		},
		{
			"From no url",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{""},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.extractSeasonNumber(tt.args.in); got != tt.want {
				t.Errorf("TmdbMIProvider.extractSeasonNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_extractSeasonTotalEpisodes(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		in string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"From actual content",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{`

                    2008 • 7 Episodes
                  `},
			7,
		},
		{
			"From clean content",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{`2099 43 Episodes`},
			43,
		},
		{
			"From no content",
			fields{"tmdb", scrapper.NewMockGoQueryScrapper()},
			args{""},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.extractSeasonTotalEpisodes(tt.args.in); got != tt.want {
				t.Errorf("TmdbMIProvider.extractSeasonTotalEpisodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_getSearchString(t *testing.T) {
	type args struct {
		name string
		year int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"With year",
			args{"Test Name", 2034},
			"Test+Name%20y:2034",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewMockTmdbMIProvider()
			if got := tr.getSearchString(tt.args.name, tt.args.year); got != tt.want {
				t.Errorf("TmdbMIProvider.getSearchString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmdbMIProvider_GetJellyfinCompatibleDirectoryName(t *testing.T) {
	type fields struct {
		baseUrl  string
		scrapper scrapper.Scrapper
	}
	type args struct {
		info MediaInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Simple name",
			fields: fields{
				baseUrl:  "",
				scrapper: nil,
			},
			args: args{
				info: MediaInfo{
					Name:          "Test Media",
					Description:   "",
					YearOfRelease: 2025,
					ThumbnailURL:  "",
					MediaID:       "6789",
				},
			},
			want: "Test Media (2025) [tmdbid-6789]",
		},
		{
			name: "Without year",
			fields: fields{
				baseUrl:  "",
				scrapper: nil,
			},
			args: args{
				info: MediaInfo{
					Name:          "Movie 2",
					Description:   "",
					YearOfRelease: 0,
					ThumbnailURL:  "",
					MediaID:       "1232",
				},
			},
			want: "Movie 2 [tmdbid-1232]",
		},
		{
			name: "Without media id",
			fields: fields{
				baseUrl:  "",
				scrapper: nil,
			},
			args: args{
				info: MediaInfo{
					Name:          "TV X",
					Description:   "",
					YearOfRelease: 2024,
					ThumbnailURL:  "",
					MediaID:       "",
				},
			},
			want: "TV X (2024)",
		},
		{
			name: "Without any optional info",
			fields: fields{
				baseUrl:  "",
				scrapper: nil,
			},
			args: args{
				info: MediaInfo{
					Name:          "Why though",
					Description:   "",
					YearOfRelease: -1,
					ThumbnailURL:  "",
					MediaID:       "",
				},
			},
			want: "Why though",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TmdbMIProvider{
				baseUrl:  tt.fields.baseUrl,
				scrapper: tt.fields.scrapper,
			}
			if got := tr.GetJellyfinCompatibleDirectoryName(tt.args.info); got != tt.want {
				t.Errorf("TmdbMIProvider.GetJellyfinCompatibleDirectoryName() = %v, want %v", got, tt.want)
			}
		})
	}
}
