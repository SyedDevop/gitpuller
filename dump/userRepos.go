package main

import (
	"encoding/json"
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

	if err := util.CreateDir(basePath); err != nil {
		return nil, err
	}

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(body)
	if err != nil {
		return nil, err
	}

	var repos []git.Repos
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
