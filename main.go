package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/api"
	"github.com/MilindGour/jellyfin-media-renamer/middlewares"
	"github.com/MilindGour/jellyfin-media-renamer/util"
	"github.com/gorilla/mux"
)

func main() {
	applicationPort := ":7749"
	log.SetFlags(0)

	// get command line args
	args := os.Args
	for _, arg := range args {
		if arg == "--dev" {
			util.SetEnvironment(util.DEV)
		}
	}

	// some informational logs
	if util.IsProduction() {
		// Production mode
		log.Println("JMR running in production mode")
	} else {
		// Dev mode
		log.Println("JMR running in dev mode")
	}
	log.Println("Using config:", util.GetConfigFilename())
	log.Println("Server starting on", applicationPort)

	r := mux.NewRouter()
	r.Use(middlewares.CorsMW)

	// API subrouter
	apiSubrouter := r.PathPrefix("/api").Subrouter()
	api.RegisterAPIRoutes(apiSubrouter)

	if !util.IsProduction() {
		corsed := middlewares.CorsMW(r)
		http.ListenAndServe(applicationPort, corsed)
	} else {
		http.ListenAndServe(applicationPort, r)
	}

}
