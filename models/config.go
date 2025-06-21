package models

type Config struct {
	Version string         `json:"version"`
	Source  []ConfigSource `json:"source"`
}

type ConfigSource struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
