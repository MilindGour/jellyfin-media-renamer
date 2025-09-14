package network

import "net/http"

type HttpHtml struct {
}

func (h *HttpHtml) GetHTML(url string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Accept-Language", "en-US")

	client := http.Client{}
	return client.Do(r)
}

func NewHttpHtml() *HttpHtml {
	return &HttpHtml{}
}
