package filesystem

import "fmt"

type FileSystemProvider interface {
	GetDirectorySize(DirEntry) int64
	ScanDirectory(path string, includeExtensions []string) []DirEntry
	IsMediaFile(path string) bool
	IsSubtitleFile(path string) bool
	MoveFiles(pathPairs []PathPair, progress chan []FileTransferProgress)
	CreateDirectory(dirpath string) bool
	DeleteDirectory(dirpath string) bool
}

type DirEntry struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	IsDirectory bool       `json:"isDirectory"`
	Children    []DirEntry `json:"children"`
}

type PathPair struct {
	OldPath string
	NewPath string
}

type FileTransferProgress struct {
	BytesTransferred int64
	PercentComplete  int
	TimeRemaining    string
	TransferSpeed    string
	RawString        string
	Error            error
}

func (ftp *FileTransferProgress) ToString() string {
	out := ""
	if ftp.Error != nil {
		out = fmt.Sprintf("FileTransfer Error: %s", ftp.Error.Error())
	} else {
		out = fmt.Sprintf("FileTransferProgress: [%d%%, %d bytes @ %s, remaining %s]", ftp.PercentComplete, ftp.BytesTransferred, ftp.TransferSpeed, ftp.TimeRemaining)
	}

	return out
}
