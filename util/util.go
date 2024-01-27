package util

import (
	"fmt"
	"os"

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
		}
	}
	return newRepos
}

// Create directory if not exists.
func CreateDir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err = os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// func
