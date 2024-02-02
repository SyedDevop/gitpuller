package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	types "github.com/SyedDevop/gitpuller/mytypes"
)

type Clint struct {
	HTTPClint  *http.Client
	GitRepoUrl string
}

func NewClint() *Clint {
	return &Clint{
		HTTPClint:  &http.Client{},
		GitRepoUrl: "",
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
		var badReq types.BadReq
		if err := json.Unmarshal(body, &badReq); err != nil {
			return err
		}
		return errors.New(badReq.Message)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func (c *Clint) GetCountents() (*[]types.Content, error) {
	if c.GitRepoUrl == "" {
		return nil, errors.New("GitRepoUrl not set")
	}
	req, err := http.NewRequest("GET", c.GitRepoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")

	var contents []types.Content

	err = c.sendRequest(req, &contents)
	if err != nil {
		return nil, err
	}
	return &contents, nil
}
