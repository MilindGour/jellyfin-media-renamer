package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIProvider interface {
	Initialize(bool)
}

type APIHandlerFn func(http.ResponseWriter, *http.Request)

func ToJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Cannot marshal to json: %v", err.Error())
		return nil
	}
	return data
}
