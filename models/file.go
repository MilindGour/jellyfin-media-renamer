package models

type DirectoryEntry struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Path        string           `json:"path"`
	Size        int64            `json:"size"`
	IsDirectory bool             `json:"isDirectory"`
	Children    []DirectoryEntry `json:"children"`
}
