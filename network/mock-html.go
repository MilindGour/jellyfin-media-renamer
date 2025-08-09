package network

import (
	"net/http"
	"net/http/httptest"

	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

type MockHtml struct {
	mockHtmlMap map[string][]byte
}

func NewMockHtml() *MockHtml {
	return &MockHtml{
		mockHtmlMap: map[string][]byte{
			"mock-scrap-html": testdata.MockScrapHtml,
		},
	}
}

func (h *MockHtml) GetHTML(url string) (*http.Response, error) {
	var resBytes []byte
	if url == "mock-scrap-html" {
		resBytes = testdata.MockScrapHtml
	}

	w := httptest.NewRecorder()
	w.Write(resBytes)
	return w.Result(), nil
}
