package filesystem

type FileSystemProvider interface {
	GetDirectorySize(DirEntry) int64
	ScanDirectory(path string, includeExtensions []string) []DirEntry
	IsMediaFile(path string) bool
	IsSubtitleFile(path string) bool
}

type DirEntry struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	IsDirectory bool       `json:"isDirectory"`
	Children    []DirEntry `json:"children"`
}
