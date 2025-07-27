package api

import (
	"github.com/MilindGour/jellyfin-media-renamer/middlewares"
	"github.com/gorilla/mux"
)

// RegisterAPIRoutes adds /api routes to the given router.
func RegisterAPIRoutes(muxRouter *mux.Router) {
	// logging middleware
	muxRouter.Use(middlewares.LogMW)

	// /api/config
	configRouter := muxRouter.PathPrefix("/config").Subrouter()
	RegisterConfigRoutes(configRouter)

	// /api/scrap
	scrapRouter := muxRouter.PathPrefix("/scrap").Subrouter()
	RegisterScrapRoutes(scrapRouter)

	// ping route
	muxRouter.HandleFunc("/ping", HandlePingRequest)
}
