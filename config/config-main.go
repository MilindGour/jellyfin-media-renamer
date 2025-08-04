// Package config contains all the functions
// that deal with the config file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func GetConfig() (*models.Config, error) {
	return readConfigJSON()
}

func getConfigAllowedExtensions() (*models.AllowedExtensions, error) {
	theConfig, err := readConfigJSON()
	if err != nil {
		return nil, err
	}
	return &theConfig.AllowedExtensions, nil
}

// GetAllowedAllExtensions method returns the list of
// all the extensions allowed to be listed in frontend.
func GetAllowedAllExtensions() []string {
	allExt := GetAllowedMediaExtensions()
	allExt = append(allExt, GetAllowedSubtitleExtensions()...)

	return allExt
}
func GetAllowedMediaExtensions() []string {
	allowedExtensions, err := getConfigAllowedExtensions()
	if err != nil {
		return []string{".mp4", ".avi", ".mkv", ".m4v"}
	}
	return allowedExtensions.Media
}
func GetAllowedSubtitleExtensions() []string {
	allowedExtensions, err := getConfigAllowedExtensions()
	if err != nil {
		return []string{".srt"}
	}
	return allowedExtensions.Subtitle
}

func GetConfigSource() ([]models.ConfigSource, error) {
	theConfig, err := readConfigJSON()
	if err != nil {
		return nil, err
	}
	return theConfig.Source, nil
}

func GetConfigSourceByID(id int) (*models.ConfigSourceByIDResponse, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	result := util.Filter(cfg.Source, func(x models.ConfigSource) bool {
		return x.ID == id
	})

	if len(result) > 0 {
		dirPath := result[0].Path

		allExts := GetAllowedAllExtensions()
		res, err := util.GetDirectoryEntries(dirPath, allExts)
		if err != nil {
			return nil, err
		}

		out := &models.ConfigSourceByIDResponse{
			BasePath:         dirPath,
			DirectoryEntries: res,
		}

		return out, err
	}
	return nil, fmt.Errorf("cannot find the id: %d of the config", id)
}

// readConfigJSON function read the config.json file according to appropriate environment.
func readConfigJSON() (*models.Config, error) {
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
