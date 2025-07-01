package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gorilla/mux"
)

// RegisterScrapRoutes adds all the routes related to the scrapping apis.
func RegisterScrapRoutes(r *mux.Router) {
	r.HandleFunc("/search", postScrapSearch).Methods("POST")
	r.HandleFunc("/confirm", postScrapConfirmIds).Methods("POST")
}

func postScrapSearch(w http.ResponseWriter, r *http.Request) {
	in := models.ScrapSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(&in)

	allMovies := []string{}
	allTVs := []string{}
	for i := range in.CleanFilenameEntries {
		mediaName := in.CleanFilenameEntries[i].Name
		mt, ok := in.MediaTypes[i]
		if !ok {
			util.HandleAPIError(w, http.StatusBadRequest, "Cannot find mediaType for "+mediaName, nil)
		}
		if mt == models.MediaTypeTV {
			allTVs = append(allTVs, mediaName)
		} else {
			allMovies = append(allMovies, mediaName)
		}
	}

	log.Printf("postScrapSearch for Movies: [%s] and TVs: [%s]", strings.Join(allMovies, ", "), strings.Join(allTVs, ", "))

	if err != nil {
		util.HandleAPIError(w, http.StatusBadRequest, "Invalid parameter passed.", err)
		return
	}

	if state.LastSecondPageAPIResponse == nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot call /scrap/search before previous API calls", nil)
		return
	}

	// TODO: Once new scrap client is implemented, add a logic to select appropriate client.
	var scrapClient scrapper.Scrapper = scrapper.NewTmdbScrapper()

	for id, cfe := range in.CleanFilenameEntries {
		// find the mediaType
		mediaType, ok := in.MediaTypes[id]
		if !ok {
			util.HandleAPIError(w, http.StatusBadRequest, "Cannot find mediaType for id "+strconv.Itoa(id), nil)
			return
		}

		if mediaType == models.MediaTypeTV {
			// TV processing
			r, err := scrapClient.SearchTV(cfe)
			if err != nil {
				util.HandleAPIError(w, http.StatusInternalServerError, "Error searching tv", err)
				return
			}
			state.LastSecondPageAPIResponse.TVResults[id] = r
		} else {
			// Movie processing
			r, err := scrapClient.SearchMovie(cfe)
			if err != nil {
				util.HandleAPIError(w, http.StatusInternalServerError, "Error searching movie", err)
				return
			}
			state.LastSecondPageAPIResponse.MovieResults[id] = r
		}

	}

	responseJSON, err := json.Marshal(state.LastSecondPageAPIResponse)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal page response", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}

// postScrapConfirmIds POST /api/scrap/confirm api
// this api is used to confirm the tmdbids of the selected directories.
func postScrapConfirmIds(w http.ResponseWriter, r *http.Request) {
	// ScrapSearchConfirmRequest
	in := models.ScrapSearchConfirmRequest{}
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		util.HandleAPIError(w, http.StatusBadRequest, "Invalid argument passed.", err)
		return
	}
}
