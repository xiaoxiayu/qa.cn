package file_client

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 0
	}
	return f.Size()
}

func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}

func IsDir(filename string) bool {
	return isFileOrDir(filename, true)
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetFilelist(path_str string) []string {
	file_list := []string{}
	err := filepath.Walk(path_str, func(path_str string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		file_list = append(file_list, path_str)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return []string{}
	}
	return file_list
}
