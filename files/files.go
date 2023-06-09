package files

import (
	"os"
)

const MAX_FILE_SIZE = 25 << 20 // 25 MiB

func FileTooLarge(file string) bool {
	bytes, err := getFileSize(file)
	if err != nil {
		println(err)
		return false
	}
	return bytes > MAX_FILE_SIZE
}

func getFileSize(file string) (int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	fi, err := f.Stat()
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}
