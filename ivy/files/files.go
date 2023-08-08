package files

import (
	"bufio"
	"io"
	"os"
)

type File struct {
	Filename  string
	Filepath  string
	Diff      bool
	Size      int64
	TooLarge  bool
	NotFound  bool
	raw       string
	Lang      string
	ranges    highlightRanges
	Highlight map[int]string
}

const MAX_FILE_SIZE = 1 << 20 // 1 MiB

func FileTooLarge(file string) bool {
	bytes, err := getFileSize(file)
	if err != nil {
		return false
	}
	return bytes > MAX_FILE_SIZE
}

func getFileSize(file string) (int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}

func GetContents(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	contents := ""
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {
			contents += line
		}
		if err != nil {
			if err != io.EOF {
				println("Error (GetContents)", err)
			}
			break
		}
	}

	return contents, nil
}
