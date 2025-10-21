package filesystem

import "fmt"

type FileSystemProvider interface {
	GetDirectorySize(DirEntry) int64
	ScanDirectory(path string, includeExtensions []string) []DirEntry
	MoveFiles(pathPairs []PathPair, progress chan []FileTransferProgress)
	CreateDirectory(dirpath string) bool
	DeleteDirectory(dirpath string) bool
	GetMountPointInfo(allPaths string) MountPointInfo
}

type DirEntry struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	IsDirectory bool       `json:"isDirectory"`
	Children    []DirEntry `json:"children"`
}

type PathPair struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
}

type FileTransferProgress struct {
	BytesTransferred int64    `json:"bytes_transferred"`
	PercentComplete  int      `json:"percent_complete"`
	TimeRemaining    string   `json:"time_remaining"`
	TransferSpeed    string   `json:"transfer_speed"`
	RawString        string   `json:"raw_string"`
	Error            error    `json:"error"`
	Files            PathPair `json:"files"`
}

type MountPointInfo struct {
	MountPoint  string `json:"mount_point"`
	TotalSizeKB int64  `json:"total_size_kb"`
	FreeSizeKB  int64  `json:"free_size_kb"`
	UsedSizeKB  int64  `json:"used_size_kb"`
}

func (ftp *FileTransferProgress) ToString() string {
	out := ""
	if ftp.Error != nil {
		out = fmt.Sprintf("FileTransfer Error: %s", ftp.Error.Error())
	} else {
		out = fmt.Sprintf("FileTransferProgress: [%d%%, %d bytes @ %s, remaining %s (FROM: %s ::: TO: %s)]", ftp.PercentComplete, ftp.BytesTransferred, ftp.TransferSpeed, ftp.TimeRemaining, ftp.Files.OldPath, ftp.Files.NewPath)
	}

	return out
}
