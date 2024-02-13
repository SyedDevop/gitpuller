package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	types "github.com/SyedDevop/gitpuller/mytypes"
)

func ParseContentsUrl(path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents", path)
}

func GetRepoFromContent(contents []types.Content) []types.Repo {
	newRepos := make([]types.Repo, len(contents))
	for i, content := range contents {
		newRepos[i] = types.Repo{
			Name:        content.Name,
			Path:        content.Path,
			Size:        content.Size,
			DownloadURL: content.DownloadURL,
			Type:        content.Type,
			URL:         content.URL,
		}
	}
	return newRepos
}

// Create directory if not exists.
func CreateDir(name string) error {
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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
