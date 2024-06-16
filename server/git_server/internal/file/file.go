package file

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	JsonDataType = []map[string]interface{}
	RepoDataType = map[string]interface{}
)

func GetCurDir() (string, bool) {
	_, filename, _, ok := runtime.Caller(0)

	return filepath.Dir(filename), ok
}

func ReadData(path string) ([]byte, error) {
	fileJson := fmt.Sprintf("%s.json", path)
	curDir, ok := GetCurDir()
	if !ok {
		fmt.Println("File#ReadJson (curDir): unable to determine current file path")
		return nil, fmt.Errorf("unable to determine current file path")
	}

	file, err := os.ReadFile(filepath.Join(curDir, fileJson))
	if err != nil {
		fmt.Printf("File#ReadJson (ReadFile): %v\n", err)
		return nil, err
	}
	return file, nil
}

func ReadJson(name string, data interface{}) error {
	file, err := ReadData(name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, data)
	if err != nil {
		return err
	}
	return nil
}

// FileExist
// return True if the File Exist false if don't Exist
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetParentPath takes a file or directory path as input and attempts to extract the parent path.
//
// Parameters:
// - path: A string representing the file or directory path from which the parent path is to be extracted.
//
// Returns:
//  1. A boolean value indicating if it is a root path. True if the path is a root path, false otherwise.
//  2. A string containing the parent path if found. If a "/" is present, this will be the path up to the last "/", excluding the "/"
//     itself. If no "/" is found, indicating no parent path can be extracted, the function returns the original path.
//
// Note: This function is designed to work with UNIX-like file system paths that use "/" as a directory separator. It does not
// handle Windows paths that use "\" as a directory separator.
func GetParentPath(path string) (bool, string) {
	pathSeparator := "/"
	if runtime.GOOS == "windows" {
		pathSeparator = "\\"
	}

	index := strings.LastIndex(path, pathSeparator)
	if index == 0 || index == -1 {
		return true, path
	} else if index == len(path)-1 {
		return GetParentPath(path[:index])
	}
	return false, path[:index]
}

// /test   0
// test/   0
// test    0
// test.go 0
// test/g.go 1
// test/game 1
func GetFileDepth(path string) int {
	pathSeparator := "/"
	count := strings.Count(path, pathSeparator)
	if count == 1 &&
		(path[0] == '/' || path[len(path)-1] == '/') ||
		count == 0 {
		return 0
	}
	return count
}
