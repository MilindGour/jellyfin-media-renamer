package network

import (
	"bytes"
	"encoding/json"
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

func (h *HttpResponse) PostJSON(url string, body any, headers *http.Header) (*http.Response, error) {
	var res *http.Response

	jsonEnc, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshalling object to json. More info: %s", err.Error())
		return nil, err
	}

	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonEnc))

	if headers != nil {
		r.Header = *headers
	}
	r.Header.Add("Content-Type", "application/json")

	res, err = h.client.Do(r)

	return res, err
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		retries: 5,
		client:  &http.Client{},
	}
}
