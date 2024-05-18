package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ParseContentsUrl takes a path and a sha or branch and returns a url for git Tree.
//
// Parameters:
//  1. repoIdtenty: {owner}/{repo} Repo owner and Repo name.
//  2. sha: The sha or branch of the repo.
//
// Returns:
// - A string containing the url for git Tree.
func ParseContentsUrl(repoIdentity, sha string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/git/trees/%s", repoIdentity, sha)
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

// getGitToken retrieves a Git token from either a configuration file or environment variable.
//
// It first attempts to fetch the token from a configuration file using Viper.
//
// If the token is not found in the configuration file, it then checks for the token in an environment variable.
//
// If the token is found in neither location, an empty string is returned.
func GetGitToken() string {
	token := ""
	vToken := viper.GetString("token")
	gitToken := os.Getenv("GIT_TOKEN")
	if vToken != "" {
		token = vToken
	} else if gitToken != "" {
		token = gitToken
	}
	return token
}

func time12(t time.Time) string {
	return t.Format("Monday, 02-Jan-06 03:04:05.000 PM MST")
}
