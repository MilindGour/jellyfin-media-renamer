package models

type Config struct {
	Version string         `json:"version"`
	Source  []ConfigSource `json:"source"`
}

type ConfigSource struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
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
}
