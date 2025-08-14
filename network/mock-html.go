package network

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

type MockHtml struct {
	mockHtmlMap map[string][]byte
}

func NewMockHtml() *MockHtml {
	return &MockHtml{
		mockHtmlMap: map[string][]byte{
			"mock-scrap-html":                      testdata.MockScrapHtml,
			"tmdb/search/movie?query=test%20movie": testdata.MockTmdbMovieSearch,
			"tmdb/search/tv?query=Test%20TV":       testdata.MockTmdbTVShowSearch,
			"tmdb/tv/1396/seasons":                 testdata.MockTmdbTVShowSeasons,
			"not-found":                            []byte(""),
		},
	}
}

func (h *MockHtml) GetHTML(url string) (*http.Response, error) {
	resBytes, hasMock := h.mockHtmlMap[url]
	if !hasMock {
		resBytes, _ = h.mockHtmlMap["not-found"]
		fmt.Println("No mock found for url:", url)
	}
	if hasMock {
		fmt.Println("Mock Available URL:", url)
	}

	w := httptest.NewRecorder()
	w.Header().Add("Content-Length", strconv.Itoa(len(resBytes)))
	w.Write(resBytes)
	return w.Result(), nil
}
