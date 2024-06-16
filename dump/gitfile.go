package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/charmbracelet/log"
)

func getCurDir() (string, bool) {
	_, filename, _, ok := runtime.Caller(0)
	return filepath.Dir(filename), ok
}

type JsonDataType = []map[string]any

func getGitFile(path, toFile string) []error {
	c := getGitClient()

	errList := make([]error, 0)

	log.Info("Getting repo tree url from", "file", path)
	file, err := os.ReadFile(path)
	if err != nil {
		errList = append(errList, err)
		return errList
	}
	var datas JsonDataType
	err = json.Unmarshal(file, &datas)
	if err != nil {
		errList = append(errList, err)
		return errList
	}
	dataLen := len(datas)
	urls := make([][]string, dataLen)
	for i, da := range datas {

		if dataurlStr, ok := da["name"].(string); ok {
			urls[i] = append(urls[i], dataurlStr)
		}
		if dataStr, ok := da["trees_url"].(string); ok {
			url := fmt.Sprintf("%s/main?recursive=1", dataStr[:len(dataStr)-6])
			urls[i] = append(urls[i], url)
		}
	}

	log.Info("Done#Getting repo tree url from", "file", path)
	toFilePath := filepath.Dir(path)

	if err := util.CreateDir(filepath.Join(toFilePath, toFile)); err != nil {
		errList = append(errList, err)
		return errList
	}

	ws := sync.WaitGroup{}
	ws.Add(len(urls))
	errChan := make(chan error)

	log.Info("Fetch#Repo from GitHub")
	for _, data := range urls {
		go func(name, url string) {
			defer ws.Done()
			path := filepath.Join(toFilePath, toFile, fmt.Sprintf("%s.json", data[0]))
			file, err := os.Create(path)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()
			res, err := c.Get(data[1])
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
		}(data[0], data[1])
	}
	ws.Wait()
	close(errChan)
	for errs := range errChan {
		errList = append(errList, errs)
	}
	log.Info("Done#Repo from GitHub")
	return errList
}
