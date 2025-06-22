package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func GetConfig() (*models.Config, error) {
	return readConfigJson()
}
func GetConfigSource() ([]models.ConfigSource, error) {
	theConfig, err := readConfigJson()
	if err != nil {
		return nil, err
	}
	return theConfig.Source, nil
}

func readConfigJson() (*models.Config, error) {
	theConfig := models.Config{}
	configFileContents, err := util.GetConfigFileContents()
	if err != nil {
		return nil, errors.New(fmt.Sprint("Cannot read the contents of config.json. ", err))
	}
	err = json.Unmarshal(configFileContents, &theConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Cannot unmarshal the contents of config.json. ", err))
	}

	return &theConfig, nil
}
