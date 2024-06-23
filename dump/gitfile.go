package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/SyedDevop/gitpuller/pkg/git"
	"github.com/charmbracelet/log"
)

func getGitFile(c *client.Client, repos []git.Repos) ([]git.Tree, []error) {
	errList := make([]error, 0)
	dataLen := len(repos)
	repoPath := filepath.Join(basePath, "repo")
	treeData := make([]git.Tree, dataLen)

	if err := util.CreateDir(repoPath); err != nil {
		errList = append(errList, err)
		return nil, errList
	}

	var ws sync.WaitGroup
	errChan := make(chan error, dataLen)
	treeDataChan := make(chan git.Tree, dataLen)

	log.Info("Fetch#Repo from GitHub")
	for i, data := range repos {
		ws.Add(1)
		go func(i int, data git.Repos) {
			defer ws.Done()
			path := filepath.Join(repoPath, fmt.Sprintf("%s.json", data.Name))
			file, err := os.Create(path)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			url := data.TreesURL[:len(data.TreesURL)-6] + "/main?recursive=1"
			res, err := c.Get(url)
			if err != nil {
				errChan <- err
				return
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				errChan <- err
				return
			}
			_, err = file.Write(body)
			if err != nil {
				errChan <- err
				return
			}
			var tdata git.Tree
			err = json.Unmarshal(body, &tdata)
			if err != nil {
				errChan <- err
				return
			}
			treeDataChan <- tdata
		}(i, data)
	}

	go func() {
		ws.Wait()
		close(errChan)
		close(treeDataChan)
	}()

	for tdata := range treeDataChan {
		treeData = append(treeData, tdata)
	}

	for errs := range errChan {
		errList = append(errList, errs)
	}

	log.Info("Done#Repo from GitHub")
	return treeData, errList
}
