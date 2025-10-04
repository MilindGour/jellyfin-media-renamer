package filesystem

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
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

func (j *JmrFS) moveSingleFile(fromPath, toPath string, channel chan FileTransferProgress) {

	// start file transfer using rsync
	// TODO: add --remove-source-files in the arguments before deploy
	rsyncCmd := exec.Command("rsync", "-avz", "--info=progress2", fromPath, toPath)
	stdOutPipe, errOut := rsyncCmd.StdoutPipe()
	stdErrPipe, errErr := rsyncCmd.StderrPipe()

	defer stdOutPipe.Close()
	defer stdErrPipe.Close()

	if errOut != nil {
		channel <- FileTransferProgress{
			Error: errOut,
		}
		close(channel)
		return
	}
	if errErr != nil {
		channel <- FileTransferProgress{
			Error: errErr,
		}
		close(channel)
		return
	}

	if errStart := rsyncCmd.Start(); errStart != nil {
		channel <- FileTransferProgress{
			Error: errStart,
		}
		close(channel)
		return
	}

	// StdOut scanner
	outScanner := bufio.NewScanner(stdOutPipe)
	outScanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := range data {
			if data[i] == '\r' || data[i] == '\n' {
				return i + 1, data[:i], nil
			}
		}

		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil

	})

	for outScanner.Scan() {
		line := outScanner.Text()

		progress := j.parseRsyncOutputToProgress(line)
		if progress != nil {
			channel <- *progress
		}
	}

	// StdErr scanner
	errScanner := bufio.NewScanner(stdErrPipe)
	for errScanner.Scan() {
		line := errScanner.Text()

		channel <- FileTransferProgress{
			Error: errors.New(line),
		}
	}

	// Wait for rsync command to finish
	if errWait := rsyncCmd.Wait(); errWait != nil {
		channel <- FileTransferProgress{
			Error: errWait,
		}
		close(channel)
		return
	}

	channel <- FileTransferProgress{
		PercentComplete: 100,
	}
	close(channel)
}

func (j *JmrFS) MoveFiles(pathPairs []PathPair, progressChannel chan []FileTransferProgress) {
	ftp := make([]FileTransferProgress, len(pathPairs))
	for pIndex, p := range pathPairs {
		newDir := path.Dir(p.NewPath)

		if j.CreateDirectory(newDir) == true {
			ch := make(chan FileTransferProgress)
			go j.moveSingleFile(p.OldPath, p.NewPath, ch)
			for progress := range ch {
				ftp[pIndex] = progress
				progressChannel <- ftp
			}
		} else {
			log.Printf("Failed to create directory %s. Skipping file %s", newDir, p.OldPath)
		}
	}
	close(progressChannel)
}

func (j *JmrFS) CreateDirectory(dirpath string) bool {
	mkdirCmd := exec.Command("mkdir", "-p", dirpath)
	err := mkdirCmd.Run()

	if err != nil {
		log.Printf("Error creating directory %s. More info:\n %s", dirpath, err.Error())
		return false
	}
	return true
}

func (j *JmrFS) DeleteDirectory(dirpath string) bool {
	rmCmd := exec.Command("rm", "-rf", dirpath)
	err := rmCmd.Run()

	if err != nil {
		log.Printf("Error deleting file / directory %s. More info:\n%s", dirpath, err.Error())
		return false
	}
	return true
}

func (j *JmrFS) parseRsyncOutputToProgress(outputLine string) *FileTransferProgress {
	splits := strings.Split(outputLine, " ")
	requiredSplits := []string{}

	for _, split := range splits {
		t := strings.Trim(split, " \t")
		if len(t) > 0 {
			requiredSplits = append(requiredSplits, t)
		}
	}

	testRE := regexp.MustCompile(`^[\d,]+$`)
	if len(requiredSplits) == 4 && testRE.MatchString(requiredSplits[0]) {
		// This IS a progress line
		out := FileTransferProgress{}
		// extract bytesTransferred
		bts := strings.ReplaceAll(requiredSplits[0], ",", "")
		var err error
		out.BytesTransferred, err = strconv.ParseInt(bts, 10, 64)
		if err != nil {
			return nil
		}

		// extract percentCompleted
		pcs := strings.ReplaceAll(requiredSplits[1], "%", "")
		out.PercentComplete, err = strconv.Atoi(pcs)
		if err != nil {
			return nil
		}

		// extract transferSpeed
		out.TransferSpeed = requiredSplits[2]

		// extract timeRemaining
		out.TimeRemaining = requiredSplits[3]

		return &out
	}

	return nil
}
