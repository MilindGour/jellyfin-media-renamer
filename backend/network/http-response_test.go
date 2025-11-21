package network

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetMockHTMLServer() *httptest.Server {
	retry := 1

	sm := http.NewServeMux()
	sm.HandleFunc("GET /pass-in-one-go", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	})
	sm.HandleFunc("GET /pass-in-two-go", func(w http.ResponseWriter, r *http.Request) {
		if retry < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Try again"))
			retry++
			return
		}
		retry = 0
		w.Write([]byte("Success in 2"))
	})
	sm.HandleFunc("GET /never-pass", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error text"))
	})

	return httptest.NewServer(sm)
}

func TestHttpHtml_GetHTML(t *testing.T) {
	// mock server init
	mockServer := GetMockHTMLServer()
	mockServerUrl := mockServer.URL
	log.Printf("Created mock server on url: %s", mockServerUrl)

	type fields struct {
		retries int
		client  *http.Client
	}
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "Should pass in 1 attempt",
			fields: fields{
				retries: 1,
				client:  mockServer.Client(),
			},
			args: args{
				url: fmt.Sprintf("%s/pass-in-one-go", mockServerUrl),
			},
			wantText: "Success",
			wantErr:  false,
		},
		{
			name: "Should pass in 2 attempt",
			fields: fields{
				retries: 2,
				client:  mockServer.Client(),
			},
			args: args{
				url: fmt.Sprintf("%s/pass-in-two-go", mockServerUrl),
			},
			wantText: "Success in 2",
			wantErr:  false,
		},
		{
			name: "Should fail",
			fields: fields{
				retries: 5,
				client:  mockServer.Client(),
			},
			args: args{
				url: fmt.Sprintf("%s/never-pass", mockServerUrl),
			},
			wantText: "error text",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpResponse{
				retries: tt.fields.retries,
				client:  tt.fields.client,
			}
			got, err := h.GetResponse(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpHtml.GetHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// extract text
			allBytes, err := io.ReadAll(got.Body)
			if err != nil {
				t.Errorf("HttpHtml.GetHTML() err while reading body: %s", err.Error())
				return
			}
			gotTxt := string(allBytes)
			if tt.wantText != gotTxt {
				t.Errorf("HttpHtml.GetHTML()= %v, want %v", gotTxt, tt.wantText)
				return
			}
		})
	}

	mockServer.Close()
}
