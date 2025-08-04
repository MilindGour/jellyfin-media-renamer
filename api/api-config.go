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
	r.HandleFunc("/source/{id}", getConfigSourceByID).Methods("GET")
	r.HandleFunc("/source/select", postSelectConfigSourceID).Methods("POST").Headers("Content-Type", "text/plain")
}

func getConfig(w http.ResponseWriter, _ *http.Request) {
	jmrConfig := config.NewJmrConfig()

	cfgjson, err := json.Marshal(jmrConfig)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal config", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}

// getConfigSource GET /api/config/source.
func getConfigSource(w http.ResponseWriter, r *http.Request) {
	sourceList := config.NewJmrConfig().GetSourceList()

	cfgjson, err := json.Marshal(sourceList)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal config source", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(cfgjson)
}

// getConfigSourceByID GET /api/config/source/:id
func getConfigSourceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if ok {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			util.HandleAPIError(w, http.StatusInternalServerError, "Cannot convert id to integer", err)
			return
		}
		// reset fileID to enable predictive ids
		util.ResetNextFileID()
		jmrConfig := config.NewJmrConfig()
		res, err := jmrConfig.GetSourceByID(idInt)
		if err != nil {
			state.LastConfigSourceByID = nil
			util.HandleAPIError(w, http.StatusNotFound, "Cannot find the specified ID", err)
			return
		}
		state.LastConfigSourceByID = res
		jsonResult, err := json.Marshal(res)
		if err != nil {
			util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal directory entries", err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonResult)
	} else {
		util.HandleAPIError(w, http.StatusBadRequest, "id param not found in the path.", nil)
	}
}

// postSelectConfigSourceID POST /api/config/source/select.
// This api is used to confirm the selected directory Ids by user for processing
// in the next screen (Page 2).
func postSelectConfigSourceID(w http.ResponseWriter, r *http.Request) {
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
	for _, theID := range selectedIdsArray {
		idInt, err := strconv.Atoi(theID)
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

	responseJSON, err := json.Marshal(response)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal page response", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
