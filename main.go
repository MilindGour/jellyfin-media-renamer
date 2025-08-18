package main

import (
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/app"
)

func main() {
	isDevEnv := false

	for _, arg := range os.Args {
		if arg == "--dev" {
			isDevEnv = true
		}
	}

	app := app.NewJmrApplication(isDevEnv)
	app.Run()
}
