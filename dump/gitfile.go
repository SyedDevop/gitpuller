package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/SyedDevop/gitpuller/pkg/git"
	"github.com/charmbracelet/log"
)

// func getCurDir() (string, bool) {
// 	_, filename, _, ok := runtime.Caller(0)
// 	return filepath.Dir(filename), ok
// }

func getGitFile(c *client.Client, repos []git.Repos) []error {
	errList := make([]error, 0)
	dataLen := len(repos)
	repoPath := filepath.Join(basePath, "repo")

	if err := util.CreateDir(repoPath); err != nil {
		errList = append(errList, err)
		return errList
	}

	ws := sync.WaitGroup{}
	ws.Add(dataLen)
	errChan := make(chan error)
	TreeData := make(chan git.Tree)

	log.Info("Fetch#Repo from GitHub")
	go func() {
		per := math.Floor((float64(len(TreeData)) / float64(dataLen)) * 100)
		log.Infof("\rOn Downloading repos .... |%0.f| %d/%d", per, len(TreeData), dataLen)
	}()
	for _, data := range repos {
		go func() {
			defer ws.Done()
			path := filepath.Join(repoPath, fmt.Sprintf("%s.json", data.Name))
			file, err := os.Create(path)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()
			res, err := c.Get(data.TreesURL[:len(data.TagsURL)-6])
			if err != nil {
				errChan <- err
				return
			}
			defer res.Body.Close()
			_, err = io.Copy(file, res.Body)
			if err != nil {
				errChan <- err
				return
			}
			var tdata git.Tree
			client.UnmarshalJSON(res, &tdata)
			TreeData <- tdata
		}()
	}
	ws.Wait()
	close(errChan)
	for errs := range errChan {
		errList = append(errList, errs)
	}
	log.Info("Done#Repo from GitHub")
	return errList
}
