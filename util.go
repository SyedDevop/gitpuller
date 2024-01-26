package main

import (
	"fmt"

	types "github.com/SyedDevop/gitpuller/mytypes"
)

func parseContentsUrl(path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents", path)
}

func getRepoFromContent(contents []types.Content) []types.Repo {
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
