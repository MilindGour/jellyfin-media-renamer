package app

import (
	"github.com/MilindGour/jellyfin-media-renamer/api"
	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
)

type JmrApplication struct {
	cfg config.ConfigProvider
	api api.APIProvider

	devMode bool
}

func NewJmrApplication(isDev bool) *JmrApplication {
	fsProvider := filesystem.NewJmrFS()

	if isDev {
		// DEV mode
		configProvider := config.NewDevJmrConfig()

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider),
			devMode: true,
		}
	} else {
		// PROD mode
		configProvider := config.NewJmrConfig()

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider),
			devMode: false,
		}
	}
}

func (app *JmrApplication) Run() {
	app.api.Initialize(true)
}
