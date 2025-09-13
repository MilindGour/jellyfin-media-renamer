package app

import (
	"github.com/MilindGour/jellyfin-media-renamer/api"
	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
	"github.com/MilindGour/jellyfin-media-renamer/renamer"
)

type JmrApplication struct {
	cfg config.ConfigProvider
	api api.APIProvider
	ren renamer.Renamer

	devMode bool
}

func NewJmrApplication(isDev bool) *JmrApplication {
	fsProvider := filesystem.NewJmrFS()

	if isDev {
		// DEV mode
		configProvider := config.NewDevJmrConfig()
		mediaInfoProvider := mediainfoprovider.NewMockTmdbMIProvider()
		ren := renamer.NewJmrRenamerV1(mediaInfoProvider)

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider, ren),
			ren:     ren,
			devMode: true,
		}
	} else {
		// PROD mode
		configProvider := config.NewJmrConfig()
		mediaInfoProvider := mediainfoprovider.NewTmdbMIProvider()
		ren := renamer.NewJmrRenamerV1(mediaInfoProvider)

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider, ren),
			ren:     ren,
			devMode: false,
		}
	}
}

func (app *JmrApplication) Run() {
	app.api.Initialize(true)
}
