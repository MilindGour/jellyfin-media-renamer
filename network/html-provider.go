package network

import "net/http"

type HtmlProvider interface {
	GetHTML(url string) (*http.Response, error)
}
