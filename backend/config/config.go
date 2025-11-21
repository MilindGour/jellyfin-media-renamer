// Package config contains all the functions
// that deal with the config file.
package config

import mediainfoprovider "github.com/MilindGour/jellyfin-media-renamer/mediaInfoProvider"

type ConfigProvider interface {
	GetConfig() *Config
	GetSourceList() []DirConfig
	GetDestinationList() []DestConfig
	GetPort() string
	GetAllowedExtensions() []string
	GetMediaExtensions() []string
	GetSubtitleExtensions() []string
	ParseFromBytes(config []byte) *Config
	ParseFromFilename(filename string) *Config
}

type DirConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type DestConfig struct {
	Name       string                      `json:"name"`
	Path       string                      `json:"path"`
	Type       mediainfoprovider.MediaType `json:"type"`
	ID         int                         `json:"id"`
	MountPoint string                      `json:"mount_point"`
}

type AllowedExtensions struct {
	Subtitle []string `json:"subtitle"`
	Media    []string `json:"media"`
}

type Config struct {
	Version           string            `json:"version"`
	Port              string            `json:"port"`
	AllowedExtensions AllowedExtensions `json:"allowedExtensions"`
	Source            []DirConfig       `json:"source"`
	Destination       []DestConfig      `json:"destination"`
}
