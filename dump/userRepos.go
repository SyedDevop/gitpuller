package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/SyedDevop/gitpuller/pkg/git"
)

func reposUrl(name string) string {
	per := 100
	pages := 1
	gitToken := util.GetGitToken()
	if gitToken == "" {
		return git.AddPaginationParams(git.UserReposURL(name), &per, &pages)
	}
	return git.AddPaginationParams(git.AuthReposURL(), &per, &pages)
}

func userRepos(c *client.Client) ([]git.Repos, error) {
	fileName := fmt.Sprintf("%s.json", *user)
	fileLocation := filepath.Join(basePath, fileName)

	res, err := c.Get(reposUrl(*user))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	file, err := os.Create(fileLocation)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return nil, err
	}

	var repos []git.Repos
	if err := client.UnmarshalJSON(res, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}
