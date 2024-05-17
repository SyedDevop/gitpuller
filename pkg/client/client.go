package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/git"
	tea "github.com/charmbracelet/bubbletea"
)

type Clint struct {
	HTTPClint  *http.Client
	GitRepoUrl string
	GitToken   string
}

func NewClint() *Clint {
	return &Clint{
		HTTPClint:  &http.Client{},
		GitRepoUrl: "",
		GitToken:   util.GetGitToken(),
	}
}

func time12(t time.Time) string {
	return t.Format("Monday, 02-Jan-06 03:04:05.000 PM MST")
}

func (c *Clint) sendRequest(req *http.Request, v interface{}) error {
	logF, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer logF.Close()
	startTime := time.Now()

	fmt.Fprintf(logF, "GetCountents#Requast Started: @ = %s \n", time12(startTime))
	res, err := c.HTTPClint.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Fprintf(logF, "GetCountents#Requast Completed: took = %d:ms finished @ = %s \n", time.Duration(time.Since(startTime)).Milliseconds(), time12(time.Now()))

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Fprintf(logF, "GetCountents#Requast red to buffer: took = %d:ms finished @ = %s \n", time.Duration(time.Since(startTime)).Milliseconds(), time12(time.Now()))
	if res.StatusCode != http.StatusOK {
		var badReq git.BadReq
		if err := json.Unmarshal(body, &badReq); err != nil {
			return err
		}

		errorMsg := fmt.Sprintf("request failed with status code %d and error message: '%s'. if the repository is private, please verify your access rights or temporary service outages.", res.StatusCode, badReq.Message)
		return fmt.Errorf(errorMsg)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	fmt.Fprintf(logF, "GetCountents#Requast finished marshaling: took = %d:ms finished @ = %s \n", time.Duration(time.Since(startTime)).Milliseconds(), time12(time.Now()))
	return nil
}

func (c *Clint) GetCountents(url *string) ([]git.TreeElement, error) {
	if c.GitRepoUrl == "" {
		return nil, errors.New("GitRepoUrl not set")
	}
	if url == nil {
		url = &c.GitRepoUrl
	}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	if c.GitToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.GitToken)
	}
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	var contents git.Tree

	err = c.sendRequest(req, &contents)
	if err != nil {
		return nil, err
	}
	return contents.Tree, nil
}
