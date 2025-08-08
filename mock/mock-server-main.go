// Package mock has the functions to create mock server for TMDB
package mock

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/scrapper_old"
)

type MockTMDBServer struct {
	http.Handler
	baseURL string
}

func NewMockTmdbClient() (*scrapper_old.TmdbScrapper, *httptest.Server) {
	client := scrapper_old.NewTmdbScrapper()
	mockServer := newMockTmdbServer()
	client.BaseURL = mockServer.URL

	return client, mockServer
}

func newMockTmdbServer() *httptest.Server {
	return httptest.NewServer(MockTMDBServer{
		baseURL: "http://mock.tmdb.local",
		Handler: createMockTmdbHandler(),
	})
}

func createMockTmdbHandler() http.Handler {
	ms := http.NewServeMux()

	addSearchMovieMockAPI(ms)

	return ms
}

func addSearchMovieMockAPI(mux *http.ServeMux) {
	mux.HandleFunc("/search/movie", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		log.Println("Got request for movie ", query)
		filename := ""

		switch {
		case strings.HasPrefix(query, "Airplane"):
			filename = "../testdata/tmdb_search_movie_airplane.html"
		}

		http.ServeFile(w, r, filename)
	})
}
