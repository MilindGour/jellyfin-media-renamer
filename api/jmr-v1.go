package api

import (
	"log"

	"github.com/MilindGour/jellyfin-media-renamer/config"
	"github.com/gorilla/mux"
)

type JmrAPI struct {
	// any future api keys or such will
	// appear here
	*mux.Router
	configProvider config.ConfigProvider
}

func NewJmrApi(isDevEnv bool) *JmrAPI {
	var targetCfg config.ConfigProvider

	if isDevEnv {
		log.Println("JMR API starting in DEVELOPER environment")
		targetCfg = config.NewDevJmrConfig()

	} else {
		log.Println("JMR API starting in PRODUCTION environment")
		targetCfg = config.NewJmrConfig()
	}

	jmrApi := JmrAPI{
		Router:         mux.NewRouter(),
		configProvider: targetCfg,
	}

	if jmrApi.configProvider == nil {
		log.Fatal("Please place config before running the server")
	}

	log.Printf("Starting server on port %s", jmrApi.GetPort())
	return &jmrApi
}
func (j *JmrAPI) GetPort() string {
	return j.configProvider.GetPort()
}

func (j *JmrAPI) RegisterAPIRoutes(mux *mux.Router) {
}
