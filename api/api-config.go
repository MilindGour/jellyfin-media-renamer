// Package api contains all the APIs configured for the project.
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gorilla/mux"
)

// RegisterConfigRoutes adds all the routes related to the config apis.
func RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("", getConfig).Methods("GET")
	r.HandleFunc("/source", getConfigSource).Methods("GET")
	r.HandleFunc("/source/{id}", getConfigSourceById).Methods("GET")
	r.HandleFunc("/source/select", postSelectConfigSourceId).Methods("POST").Headers("Content-Type", "text/plain")
}

func getConfig(w http.ResponseWriter, _ *http.Request) {
	cfg, err := config.GetConfig()
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot get config", err)
		return
	}

	cfgjson, err := json.Marshal(cfg)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal config", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}

// getConfigSource GET /api/config/source.
func getConfigSource(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetConfigSource()
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot get config source", err)
		return
	}

	cfgjson, err := json.Marshal(cfg)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal config source", err)
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
			util.HandleAPIError(w, http.StatusInternalServerError, "Cannot convert id to integer", err)
			return
		}
		dirEntries, err := config.GetConfigSourceById(idInt)
		if err != nil {
			util.HandleAPIError(w, http.StatusInternalServerError, "Cannot get config by id", err)
			return
		}
		jsonResult, err := json.Marshal(dirEntries)
		if err != nil {
			util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal directory entries", err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonResult)
	} else {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot read id param from url", nil)
	}
}

// postSelectConfigSourceId POST /api/config/source/select.
// This api is used to confirm the selected directory Ids by user for processing
// in the next screen (Page 2).
func postSelectConfigSourceId(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var selectedIds string
	_, err := fmt.Fscanf(r.Body, "%s", &selectedIds)
	if err != nil {
		util.HandleAPIError(w, http.StatusBadRequest, "Selected Ids not passed", err)
		return
	}
	log.Printf("selectedIds: %s", selectedIds)
	selectedIdsArray := strings.Split(selectedIds, ",")

	if len(selectedIdsArray) == 0 {
		util.HandleAPIError(w, http.StatusBadRequest, "Selected Ids not passed", nil)
		return
	}

	var selectedIdsInt []int

	// change to int and raise if not convertible
	for _, theId := range selectedIdsArray {
		idInt, err := strconv.Atoi(theId)
		if err != nil {
			util.HandleAPIError(w, http.StatusBadRequest, "Only comma separated integers are allowed", err)
			return
		}
		selectedIdsInt = append(selectedIdsInt, idInt)
	}
	// Do some operation and return StatusNoContent
	response, err := config.PopulateSecondScreenResponse(selectedIdsInt)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Error populating second page response", err)
	}

	// Store response so that further results can use data from it.
	state.LastSecondPageAPIResponse = response

	responseJson, err := json.Marshal(response)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal page response", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJson)
}
