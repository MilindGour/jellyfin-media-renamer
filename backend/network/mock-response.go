package network

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

type MockResponse struct {
	mockHtmlMap map[string][]byte
}

func NewMockResponse() *MockResponse {
	return &MockResponse{
		mockHtmlMap: map[string][]byte{
			"mock-scrap-html":                           testdata.MockScrapHtml,
			"tmdb/search/movie?query=test+movie":        testdata.MockTmdbMovieSearch,
			"tmdb/search/tv?query=Test+TV":              testdata.MockTmdbTVShowSearch,
			"tmdb/tv/1396/seasons":                      testdata.MockTmdbTVShowSeasons,
			"tmdb/movie/872585":                         testdata.MockTmdbMovieDetail,
			"https://apibay.org/q.php?q=TestSearchTerm": testdata.MockNewMediaSearchResponse,
			"not-found":                                 []byte(""),
		},
	}
}

func (h *MockResponse) GetResponse(url string) (*http.Response, error) {
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

func (h *MockResponse) PostJSON(url string, body any, headers *http.Header) (*http.Response, error) {
	return nil, nil
}
