package testdata

import (
	"embed"
	_ "embed"
)

//go:embed tmdb_scrap_search_response.json
var ScrapSearchResponseMock []byte

//go:embed tmdb_rename_confirm_request.json
var RenameConfirmRequestMock []byte

//go:embed config.test.json
var ConfigJsonMock []byte

// New refactored things go below

//go:embed html/mock-scrap-html.html
var MockScrapHtml []byte

//go:embed html/mock-tmdb-movie-search.html
var MockTmdbMovieSearch []byte

//go:embed html/mock-tmdb-tvshow-search.html
var MockTmdbTVShowSearch []byte

//go:embed html/mock-tmdb-tvshow-seasons.html
var MockTmdbTVShowSeasons []byte

//go:embed fs-structure/*
var MockFSStructure embed.FS

//go:embed html/mock-tmdb-movie-detail.html
var MockTmdbMovieDetail []byte
