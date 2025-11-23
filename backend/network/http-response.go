package network

import (
	"fmt"
	"net/http"
)

type HttpResponse struct {
	retries int
	client  *http.Client
}

func (h *HttpResponse) GetResponse(url string) (*http.Response, error) {
	var err error
	var res *http.Response

	for range h.retries {
		r, _ := http.NewRequest(http.MethodGet, url, nil)
		r.Header.Add("Accept-Language", "en-US")
		res, err = h.client.Do(r)

		if err == nil && res.StatusCode == http.StatusOK {
			return res, err
		}
	}

	return res, fmt.Errorf("Cannot get html from url: %s, even after %d attempts", url, h.retries)
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		retries: 5,
		client:  &http.Client{},
	}
}
