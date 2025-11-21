package api

import "net/http"

func HandlePingRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("online"))
}
