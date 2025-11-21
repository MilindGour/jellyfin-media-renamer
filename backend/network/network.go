package network

import "net/http"

type HttpResponseProvider interface {
	GetResponse(url string) (*http.Response, error)
}
