package network

import "net/http"

type HttpHtml struct {
}

func (h *HttpHtml) GetHTML(url string) (*http.Response, error) {
	return http.Get(url)
}
