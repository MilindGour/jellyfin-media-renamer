package models

type DirectoryEntry struct {
	Name        string           `json:"name"`
	Size        int64            `json:"size"`
	IsDirectory bool             `json:"isDirectory"`
	Children    []DirectoryEntry `json:"children"`
}
