package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/state"
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

func GetConfigSourceById(id int) ([]models.DirectoryEntry, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	result := util.Filter(cfg.Source, func(x models.ConfigSource) bool {
		return x.Id == id
	})

	if len(result) > 0 {
		dirPath := result[0].Path
		// Storing in LastConfigSourceById for next api call
		state.LastConfigSourceById, err = util.GetDirectoryEntries(dirPath)
		return state.LastConfigSourceById, err
	}
	return nil, errors.New(fmt.Sprintf("Cannot find the id: %d of the config", id))
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
