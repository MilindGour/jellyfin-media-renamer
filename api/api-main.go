package api

import (
	"github.com/gorilla/mux"
)

// RegisterAPIRoutes adds /api routes to the given router.
func RegisterAPIRoutes(muxRouter *mux.Router) {

	// /api/config
	configRouter := muxRouter.PathPrefix("/config").Subrouter()
	RegisterConfigRoutes(configRouter)
}
