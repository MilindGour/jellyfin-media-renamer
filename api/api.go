package api

import "github.com/gorilla/mux"

type APIProvider interface {
	RegisterAPIRoutes(mux *mux.Router)
	GetPort() string
}
