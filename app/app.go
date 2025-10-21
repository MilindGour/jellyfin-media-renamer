package app

import (
	"github.com/MilindGour/jellyfin-media-renamer/api"
	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"
	"github.com/MilindGour/jellyfin-media-renamer/renamer"
	"github.com/MilindGour/jellyfin-media-renamer/websocket"
)

type JmrApplication struct {
	cfg config.ConfigProvider
	api api.APIProvider
	ren renamer.Renamer
	mip mediainfoprovider.MediaInfoProvider

	devMode bool
}

func NewJmrApplication(isDev bool) *JmrApplication {
	fsProvider := filesystem.NewJmrFS()
	ws := websocket.NewJMRWebSocket()

	if isDev {
		// DEV mode
		configProvider := config.NewDevJmrConfig()
		mediaInfoProvider := mediainfoprovider.NewTmdbMIProvider()
		ren := renamer.NewJmrRenamerV1(mediaInfoProvider, fsProvider, configProvider, ws)

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider, ren, mediaInfoProvider, ws),
			ren:     ren,
			devMode: true,
		}
	} else {
		// PROD mode
		configProvider := config.NewJmrConfig()
		mediaInfoProvider := mediainfoprovider.NewTmdbMIProvider()
		ren := renamer.NewJmrRenamerV1(mediaInfoProvider, fsProvider, configProvider, ws)

		return &JmrApplication{
			cfg:     configProvider,
			api:     api.NewJmrApi(configProvider, fsProvider, ren, mediaInfoProvider, ws),
			ren:     ren,
			devMode: false,
		}
	}
}

func (app *JmrApplication) Run() {
	app.api.Initialize(true)
}
