// Package config contains all the functions
// that deal with the config file.
package config

type ConfigProvider interface {
	GetSourceList() []DirConfig
	GetPort() string
	GetMediaExtensions() []string
	GetSubtitleExtensions() []string
	ParseFromBytes(config []byte) *Config
	ParseFromFilename(filename string) *Config
}

type DirConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
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
}
