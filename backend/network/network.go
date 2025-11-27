package network

import (
	"net/http"
)

type HttpResponseProvider interface {
	GetResponse(url string) (*http.Response, error)
	PostJSON(url string, body any, headers *http.Header) (*http.Response, error)
}
