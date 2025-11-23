package config

import (
	"encoding/json"
	"log"
	"os"
)

type JmrConfig struct {
	config *Config
}

func NewJmrConfig() *JmrConfig {
	// location of production config file must be in:
	// $HOME/.config/jmr/jmr.conf

	configLocation := "/config/jmr.config.json"
	log.Printf("Reading config from location %s", configLocation)
	cfg := JmrConfig{}
	cfg.config = cfg.ParseFromFilename(configLocation)

	if cfg.config == nil {
		log.Fatal("Cannot find the config file. Please configure that first before running the server")
	}
	return &cfg
}

func (j *JmrConfig) GetPort() string {
	if len(j.config.Port) == 0 {
		j.config.Port = "7749"
	}
	return j.config.Port
}

func (j *JmrConfig) GetConfig() *Config {
	return j.config
}

func (j *JmrConfig) GetSourceList() []DirConfig {
	out := []DirConfig{}
	for _, src := range j.config.Source {
		out = append(out, src)
	}
	return out
}

func (j *JmrConfig) GetDestinationList() []DestConfig {
	for i := range j.config.Destination {
		j.config.Destination[i].ID = i + 1
	}
	return j.config.Destination
}

func (j *JmrConfig) GetMediaExtensions() []string {
	if len(j.config.AllowedExtensions.Media) == 0 {
		j.config.AllowedExtensions.Media = []string{".mp4", ".mkv", ".m4v", ".avi"}
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
	out := j.GetMediaExtensions()
	out = append(out, j.GetSubtitleExtensions()...)
	return out
}
func (j *JmrConfig) ParseFromBytes(config []byte) *Config {
	out := &Config{}
	err := json.Unmarshal(config, out)

	if err != nil {
		return nil
	}

	return out
}
func (j *JmrConfig) ParseFromFilename(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error parsing config file. More info: %s", err.Error())
		return nil
	}
	return j.ParseFromBytes(data)
}
