package util

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
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

// func GetRepoFromContent(contents []api.Content) []api.Repo {
// 	newRepos := make([]api.Repo, len(contents))
// 	for i, content := range contents {
// 		newRepos[i] = api.Repo{
// 			Name:        content.Name,
// 			Path:        content.Path,
// 			Size:        content.Size,
// 			DownloadURL: content.DownloadURL,
// 			Type:        content.Type,
// 			URL:         content.URL,
// 		}
// 	}
// 	return newRepos
// }

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

type FileMode uint32

const (
	// Empty is used as the FileMode of tree elements when comparing
	// trees in the following situations:
	//
	// - the mode of tree elements before their creation.  - the mode of
	// tree elements after their deletion.  - the mode of unmerged
	// elements when checking the index.
	//
	// Empty has no file system equivalent.  As Empty is the zero value
	// of FileMode, it is also returned by New and
	// NewFromOsNewFromOSFileMode along with an error, when they fail.
	Empty FileMode = 0
	// Dir represent a Directory.
	Dir FileMode = 0040000
	// Regular represent non-executable files.  Please note this is not
	// the same as golang regular files, which include executable files.
	Regular FileMode = 0100644
	// Deprecated represent non-executable files with the group writable
	// bit set.  This mode was supported by the first versions of git,
	// but it has been deprecated nowadays.  This library uses them
	// internally, so you can read old packfiles, but will treat them as
	// Regulars when interfacing with the outside world.  This is the
	// standard git behaviour.
	Deprecated FileMode = 0100664
	// Executable represents executable files.
	Executable FileMode = 0100755
	// Symlink represents symbolic links to files.
	Symlink FileMode = 0120000
	// Submodule represents git submodules.  This mode has no file system
	// equivalent.
	Submodule FileMode = 0160000
)

// New takes the octal string representation of a FileMode and returns
// the FileMode and a nil error.  If the string can not be parsed to a
// 32 bit unsigned octal number, it returns Empty and the parsing error.
//
// Example: "40000" means Dir, "100644" means Regular.
//
// Please note this function does not check if the returned FileMode
// is valid in git or if it is malformed.  For instance, "1" will
// return the malformed FileMode(1) and a nil error.
func FileModeNew(s string) (FileMode, error) {
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return Empty, err
	}

	return FileMode(n), nil
}
