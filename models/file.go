package models

type DirectoryEntry struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Size        int64            `json:"size"`
	IsDirectory bool             `json:"isDirectory"`
	Children    []DirectoryEntry `json:"children"`
}
