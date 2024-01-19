package main

import (
	"fmt"

	. "github.com/SyedDevop/gitpuller/mytypes"
)

func parseContentsUrl(path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents", path)
}

func getRepoFromContent(contents []Content) []Repo {
	newRepos := make([]Repo, len(contents))
	for i, content := range contents {
		newRepos[i] = Repo{
			Name:        content.Name,
			Path:        content.Path,
			Size:        content.Size,
			DownloadURL: content.DownloadURL,
			Type:        content.Type,
		}
	}
	return newRepos
}
