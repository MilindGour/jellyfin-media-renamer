package main

import (
	"github.com/MilindGour/jellyfin-media-renamer/api"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// API subrouter
	apiSubrouter := r.PathPrefix("/api").Subrouter()
	api.RegisterAPIRoutes(apiSubrouter)

	http.ListenAndServe(":7749", r)
}
