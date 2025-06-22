// Package api contains all the APIs configured for the project.
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gorilla/mux"
)

// RegisterConfigRoutes adds all the routes related to the config apis.
func RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("", getConfig)
	r.HandleFunc("/source", getConfigSource)
	r.HandleFunc("/source/{id}", getConfigSourceById)
}

func getConfig(w http.ResponseWriter, _ *http.Request) {
	cfg, err := config.GetConfig()
	if err != nil {
		util.HandleAPIError(w, 500, "Cannot get config", err)
		return
	}

	cfgjson, err := json.Marshal(cfg)
	if err != nil {
		util.HandleAPIError(w, 500, "Cannot marshal config", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}

// getConfigSource GET /api/config/source.
func getConfigSource(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetConfigSource()
	if err != nil {
		util.HandleAPIError(w, 500, "Cannot get config source", err)
		return
	}

	cfgjson, err := json.Marshal(cfg)
	if err != nil {
		util.HandleAPIError(w, 500, "Cannot marshal config source", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}

// getConfigSourceById GET /api/config/source/:id
func getConfigSourceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			util.HandleAPIError(w, 500, "Cannot convert id to integer", err)
			return
		}
		dirEntries, err := config.GetConfigSourceById(idInt)
		if err != nil {
			util.HandleAPIError(w, 500, "Cannot get config by id", err)
			return
		}
		jsonResult, err := json.Marshal(dirEntries)
		if err != nil {
			util.HandleAPIError(w, 500, "Cannot marshal directory entries", err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonResult)
	} else {
		util.HandleAPIError(w, 500, "Cannot read id param from url", nil)
	}
}
