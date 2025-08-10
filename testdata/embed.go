package testdata

import (
	_ "embed"
)

//go:embed tmdb_scrap_search_response.json
var ScrapSearchResponseMock []byte

//go:embed config.test.json
var ConfigJsonMock []byte

// New refactored things go below

//go:embed html/mock-scrap-html.html
var MockScrapHtml []byte

//go:embed html/mock-tmdb-movie-search.html
var MockTmdbMovieSearch []byte
