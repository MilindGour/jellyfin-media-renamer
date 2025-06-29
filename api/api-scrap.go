package api

import (
	"encoding/json"
	"net/http"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/scrapper"
	"github.com/MilindGour/jellyfin-media-renamer/state"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gorilla/mux"
)

// RegisterScrapRoutes adds all the routes related to the scrapping apis.
func RegisterScrapRoutes(r *mux.Router) {
	r.HandleFunc("/search", postScrapSearch).Methods("POST")
}

func postScrapSearch(w http.ResponseWriter, r *http.Request) {
	in := map[int]models.ClearFileEntry{}
	json.NewDecoder(r.Body).Decode(&in)

	if state.LastSecondPageAPIResponse == nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot call /scrap/search before previous API calls", nil)
		return
	}

	// everything should be valid at this point.
	tmdbClient := scrapper.NewTmdbScrapper()

	for id, cfe := range in {
		r, err := tmdbClient.SearchMovie(cfe)
		if err != nil {
			util.HandleAPIError(w, http.StatusInternalServerError, "Error searching movie", err)
			return
		}
		state.LastSecondPageAPIResponse.MovieResults[id] = r
	}

	responseJson, err := json.Marshal(state.LastSecondPageAPIResponse)
	if err != nil {
		util.HandleAPIError(w, http.StatusInternalServerError, "Cannot marshal page response", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJson)
}
