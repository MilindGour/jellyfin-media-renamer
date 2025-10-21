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

func (j *JmrFS) moveSingleFile(fromPath, toPath string, channel chan FileTransferProgress) {

	// start file transfer using rsync
	// TODO: add --remove-source-files in the arguments before deploy
	// and remove --bwlimit=4000
	rsyncCmd := exec.Command("rsync", "-avz", "--info=progress2", "--bwlimit=5M", fromPath, toPath)
	stdOutPipe, errOut := rsyncCmd.StdoutPipe()
	stdErrPipe, errErr := rsyncCmd.StderrPipe()

	files := PathPair{
		OldPath: fromPath,
		NewPath: toPath,
	}

	defer stdOutPipe.Close()
	defer stdErrPipe.Close()

	if errOut != nil {
		channel <- FileTransferProgress{
			Error: errOut,
			Files: files,
		}
		close(channel)
		return
	}
	if errErr != nil {
		channel <- FileTransferProgress{
			Error: errErr,
			Files: files,
		}
		close(channel)
		return
	}

	if errStart := rsyncCmd.Start(); errStart != nil {
		channel <- FileTransferProgress{
			Error: errStart,
			Files: files,
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
			progress.Files = files
			channel <- *progress
		}
	}

	// StdErr scanner
	errScanner := bufio.NewScanner(stdErrPipe)
	for errScanner.Scan() {
		line := errScanner.Text()

		channel <- FileTransferProgress{
			Error: errors.New(line),
			Files: files,
		}
	}

	// Wait for rsync command to finish
	if errWait := rsyncCmd.Wait(); errWait != nil {
		channel <- FileTransferProgress{
			Error: errWait,
			Files: files,
		}
		close(channel)
		return
	}

	// get the file length transferred
	destStat, err := os.Stat(toPath)
	if err != nil {
		channel <- FileTransferProgress{
			Error: err,
			Files: files,
		}
	}

	channel <- FileTransferProgress{
		BytesTransferred: destStat.Size(),
		PercentComplete:  100,
		Files:            files,
	}
	close(channel)
}

func (j *JmrFS) MoveFiles(pathPairs []PathPair, progressChannel chan []FileTransferProgress) {
	ftp := make([]FileTransferProgress, len(pathPairs))

	for i := range ftp {
		ftp[i].Files = pathPairs[i]
	}

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

func (j *JmrFS) GetMountPointInfo(mountPoint string) MountPointInfo {
	cmd := exec.Command("df", "-k", mountPoint)
	out := MountPointInfo{}

	stderr, er1 := cmd.StderrPipe()
	stdout, er2 := cmd.StdoutPipe()

	defer stderr.Close()
	defer stdout.Close()

	if er1 != nil {
		log.Printf("Error during StdinPipe(): %s", er1.Error())
		return out
	}
	if er2 != nil {
		log.Printf("Error during StdoutPipe(): %s", er2.Error())
		return out
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error occured while executing df -k. More info: %s", err.Error())
		return out
	}

	// StdOut scanner
	outScanner := bufio.NewScanner(stdout)
	lineCount := 0
	for outScanner.Scan() {
		line := outScanner.Text()
		lineCount++
		if lineCount == 2 {
			trimmed := strings.Trim(line, " \t\n\r")
			log.Printf("Output from df -k: %s", trimmed)
			parsedMPInfo := j.parseDFOutputToMountPointInfo(trimmed)
			if parsedMPInfo != nil {
				parsedMPInfo.MountPoint = mountPoint
				out = *parsedMPInfo
			}
		}
	}

	// StdErr scanner
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		line := outScanner.Text()
		log.Printf("Error from df -k: %s", line)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Error occured while waiting for command df -k. More info: %s", err.Error())
		return out
	}

	return out
}

func (j *JmrFS) parseDFOutputToMountPointInfo(in string) *MountPointInfo {
	splits := []string{}
	for s := range strings.SplitSeq(in, " ") {
		if len(s) > 0 {
			splits = append(splits, s)
		}
	}

	out := MountPointInfo{
		TotalSizeKB: 0,
		FreeSizeKB:  0,
		UsedSizeKB:  0,
	}
	if len(splits) >= 4 {
		var err error
		out.TotalSizeKB, err = strconv.ParseInt(splits[1], 10, 64)
		if err != nil {
			log.Printf("Error occured while parsing total size: %s", err.Error())
			return nil
		}
		out.UsedSizeKB, err = strconv.ParseInt(splits[2], 10, 64)
		if err != nil {
			log.Printf("Error occured while parsing used size: %s", err.Error())
			return nil
		}
		out.FreeSizeKB, err = strconv.ParseInt(splits[3], 10, 64)
		if err != nil {
			log.Printf("Error occured while parsing free size: %s", err.Error())
			return nil
		}

		return &out
	}

	return nil
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
