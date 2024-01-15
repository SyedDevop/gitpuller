package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Clint struct {
	HTTPClint *http.Client
}

func NewClint() *Clint {
	return &Clint{
		HTTPClint: &http.Client{},
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
		return errors.New(badReq.Message)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func (c *Clint) getCountents(url string) (*[]Content, error) {
	req, err := http.NewRequest("GET", parseContentsUrl(url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")

	var contents []Content

	err = c.sendRequest(req, &contents)
	if err != nil {
		return nil, err
	}
	return &contents, nil
}
