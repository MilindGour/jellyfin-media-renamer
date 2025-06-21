// Package api contains all the APIs configured for the project.
package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterConfigRoutes adds all the routes related to the config apis.
func RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("/source", getConfigSource)
}

// getConfigSource GET /api/config/source.
func getConfigSource(w http.ResponseWriter, r *http.Request) {
	res := []byte("Hello from config source")
	w.Write(res)
}
