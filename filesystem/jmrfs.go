package filesystem

import (
	"io/fs"
	"os"
	"path"
	"slices"
)

type JmrFS struct {
	fs fs.FS
}

func NewJmrFS() *JmrFS {
	return &JmrFS{
		fs: nil,
	}
}

func (j *JmrFS) ScanDirectory(dirpath string, includeExtensions []string) []DirEntry {
	entries, err := os.ReadDir(dirpath)
	// fs.ReadDir()
	if err != nil {
		return nil
	}

	out := []DirEntry{}
	for _, entry := range entries {
		outEntry := DirEntry{
			Name: entry.Name(),
			Path: path.Join(dirpath, entry.Name()),
		}
		if entry.IsDir() {
			// entry is a directory, recurse
			outEntry.IsDirectory = true
			outEntry.Children = j.ScanDirectory(outEntry.Path, includeExtensions)
			outEntry.Size = j.GetDirectorySize(outEntry)
		} else {
			// check if the file extension should be included
			ext := path.Ext(outEntry.Path)
			if !slices.ContainsFunc(includeExtensions, func(e string) bool {
				return ext == e
			}) {
				continue
			}

			outEntry.IsDirectory = false
			outEntry.Children = nil
			info, err := entry.Info()
			if err != nil {
				outEntry.Size = 0
			} else {
				outEntry.Size = info.Size()
			}
		}
		out = append(out, outEntry)
	}

	return out
}

func (j *JmrFS) GetDirectorySize(in DirEntry) int64 {
	out := int64(0)
	for _, child := range in.Children {
		out += child.Size
	}
	return out
}

func (j *JmrFS) IsMediaFile(path string) bool {
	panic("Not implemented")
}
func (j *JmrFS) IsSubtitleFile(path string) bool {
	panic("Not implemented")
}
