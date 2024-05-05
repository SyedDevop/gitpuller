package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Clint struct {
	HTTPClint  *http.Client
	GitRepoUrl string
	GitToken   string
}

func NewClint() *Clint {
	gitToken := os.Getenv("GIT_TOKEN")
	return &Clint{
		HTTPClint:  &http.Client{},
		GitRepoUrl: "",
		GitToken:   gitToken,
	}
}

func (c *Clint) sendRequest(req *http.Request, v interface{}) error {
	res, err := c.HTTPClint.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		var badReq BadReq
		if err := json.Unmarshal(body, &badReq); err != nil {
			return err
		}

		errorMsg := fmt.Sprintf("request failed with status code %d and error message: '%s'. if the repository is private, please verify your access rights or temporary service outages.", res.StatusCode, badReq.Message)
		return fmt.Errorf(errorMsg)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func (c *Clint) GetCountents(url *string) ([]TreeElement, error) {
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

	var contents Tree

	err = c.sendRequest(req, &contents)
	if err != nil {
		return nil, err
	}
	return contents.Tree, nil
}
