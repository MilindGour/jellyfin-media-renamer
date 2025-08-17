package filesystem

type FileSystem interface {
	ScanDirectory(path string) []DirEntry
	FilterScan(path string, fileFn FilterFunction) []DirEntry
}

type DirEntry struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Size        int        `json:"size"`
	IsDirectory bool       `json:"isDirectory"`
	Children    []DirEntry `json:"children"`
}

type FilterFunction func(DirEntry) bool
