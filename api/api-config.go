// Package api contains all the APIs configured for the project.
package api

import (
	"net/http"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/gorilla/mux"
)

// RegisterConfigRoutes adds all the routes related to the config apis.
func RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("/source", getConfigSource)
}

// getConfigSource GET /api/config/source.
func getConfigSource(w http.ResponseWriter, r *http.Request) {
	res := []byte(config.GetConfig())
	w.Write(res)
}
