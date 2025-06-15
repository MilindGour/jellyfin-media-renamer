package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAPIRoutes(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from api root"))
}
