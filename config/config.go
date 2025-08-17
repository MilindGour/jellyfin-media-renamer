// Package config contains all the functions
// that deal with the config file.
package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/MilindGour/jellyfin-media-renamer/models"
	"github.com/MilindGour/jellyfin-media-renamer/util"
)

type JmrConfig struct {
	config *models.Config
}

func (j *JmrConfig) ParseFromJsonBytes(data []byte) {
	err := json.Unmarshal(data, &j.config)
	if err != nil {
		panic("Unable to parse cofiguration. Error: " + err.Error())
	}
}
func (j *JmrConfig) ReadFromJsonFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Unable to read " + filename + ". Error: " + err.Error())
	}
	j.ParseFromJsonBytes(data)
}

func (j *JmrConfig) GetMediaExtensions() []string {
	if len(j.config.AllowedExtensions.Media) == 0 {
		j.config.AllowedExtensions.Media = []string{".mp4", ".avi", ".mkv", ".m4v"}
	}
	return j.config.AllowedExtensions.Media
}
func (j *JmrConfig) GetSubtitleExtensions() []string {
	if len(j.config.AllowedExtensions.Subtitle) == 0 {
		j.config.AllowedExtensions.Subtitle = []string{".srt"}
	}
	return j.config.AllowedExtensions.Subtitle
}
func (j *JmrConfig) GetAllowedExtensions() []string {
	allExtensions := j.GetMediaExtensions()
	allExtensions = append(allExtensions, j.GetSubtitleExtensions()...)

	return allExtensions
}
func (j *JmrConfig) GetSourceList() []models.ConfigSource {
	return j.config.Source
}
func (j *JmrConfig) GetSourceByID(id int) (*models.ConfigSourceByIDResponse, error) {
	result := util.Filter(j.config.Source, func(x models.ConfigSource) bool {
		return x.ID == id
	})

	if len(result) > 0 {
		dirPath := result[0].Path

		allExts := j.GetAllowedExtensions()
		res, err := util.GetDirectoryEntries(dirPath, allExts)
		if err != nil {
			return nil, errors.New("Error trying to get directory entries. " + err.Error())
		}

		out := &models.ConfigSourceByIDResponse{
			BasePath:         dirPath,
			DirectoryEntries: res,
		}

		return out, nil
	}
	panic("Cannot find the config source")
}

//	func NewJmrConfig() *JmrConfig {
//		if jmrConfigInstance == nil {
//			jmrConfigInstance = &JmrConfig{}
//		}
//		return jmrConfigInstance
//	}
func NewJmrConfigByFilename(filename string) *JmrConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Error reading config file. " + err.Error())
	}
	return NewJmrConfigByData(data)
}

func NewJmrConfigByData(data []byte) *JmrConfig {
	jConfig := &JmrConfig{}
	jConfig.ParseFromJsonBytes(data)

	return jConfig
}
