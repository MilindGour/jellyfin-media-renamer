package config

import (
	"encoding/json"
	"log"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func GetConfig() models.Config {
	return *readConfigJson()
}
func GetConfigSource() []models.ConfigSource {
	theConfig := readConfigJson()
	return theConfig.Source
}

func readConfigJson() *models.Config {
	theConfig := models.Config{}
	configFileContents, err := util.GetConfigFileContents()
	if err != nil {
		log.Fatal("Cannot read the contents of config.json")
	}
	err = json.Unmarshal(configFileContents, &theConfig)
	if err != nil {
		log.Fatal("Cannot unmarshal the contents of config.json")
	}

	return &theConfig
}
