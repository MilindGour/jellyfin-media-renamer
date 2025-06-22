// Package api contains all the APIs configured for the project.
package api

import (
	"encoding/json"
	"fmt"
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
	cfg, err := config.GetConfigSource()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Cannot get config.json. ", err)
		return
	}

	cfgjson, err := json.Marshal(cfg)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Cannot marshal config.source. "))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}
