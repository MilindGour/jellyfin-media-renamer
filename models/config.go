package models

type Config struct {
	Version           string            `json:"version"`
	AllowedExtensions AllowedExtensions `json:"allowedExtensions"`
	Source            []ConfigSource    `json:"source"`
}

type AllowedExtensions struct {
	Media    []string `json:"media"`
	Subtitle []string `json:"subtitle"`
}
type ConfigSource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type ConfigSourceByIDResponse struct {
	BasePath         string           `json:"basePath"`
	DirectoryEntries []DirectoryEntry `json:"directoryEntries"`
}

type ClearFileEntry struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

type SecondScreenResponse struct {
	Success              bool                   `json:"success"`
	SelectedDirEntries   []DirectoryEntry       `json:"selectedDirEntries"`
	CleanFilenameEntries map[int]ClearFileEntry `json:"cleanFilenameEntries"`
	MovieResults         map[int][]MovieResult  `json:"movieResults"`
	TVResults            map[int][]TVResult     `json:"tvResults"`
}
