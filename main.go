package main

import (
	"net/http"
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/api"
)

func main() {
	isDevEnv := false

	for _, arg := range os.Args {
		if arg == "--dev" {
			isDevEnv = true
		}
	}

	jmrAPI := api.NewJmrApi(isDevEnv)
	http.ListenAndServe(":"+jmrAPI.GetPort(), jmrAPI)
}
