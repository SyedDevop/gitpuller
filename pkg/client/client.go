package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Header     http.Header
}

func NewClint() *Client {
	return &Client{
		HTTPClient: &http.Client{},
		BaseURL:    "",
		Header:     http.Header{},
	}
}

// AddHeader adds a single header key-value pair to the Clint instance's headers.
//
// Parameters:
//   - key: The name of the header.
//   - value: The value of the header.
func (c *Client) AddHeader(key, value string) *Client {
	c.Header.Set(key, value)
	return c
}

// AddHeaders adds multiple header key-value pairs to the Clint instance's headers.
//
// Parameters:
//   - headers: A map of header key-value pairs.
func (c *Client) AddHeaders(headers map[string]string) *Client {
	for key, value := range headers {
		c.AddHeader(key, value)
	}
	return c
}

// AddBareAuth adds a Bearer authentication token to the Clint instance's headers.
//
// For Example:  ["Authorization": "Bearer <token>"]
//
// Parameters:
//   - auth: The Bearer authentication token.
func (c *Client) AddBareAuth(auth string) *Client {
	c.AddHeader("Authorization", "Bearer "+auth)
	return c
}

func (c *Client) UnmarshalJSON(res *http.Response, data interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, values := range c.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// func (c *Client) sendRequest(req *http.Request, v interface{}) error {
// 	res, err := c.HTTPClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
//
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return err
// 	}
//
// 	if res.StatusCode != http.StatusOK {
// 		var badReq git.BadReq
// 		if err := json.Unmarshal(body, &badReq); err != nil {
// 			return err
// 		}
//
// 		errorMsg := fmt.Sprintf("request failed with status code %d and error message: '%s'. if the repository is private, please verify your access rights or temporary service outages.", res.StatusCode, badReq.Message)
// 		return fmt.Errorf(errorMsg)
// 	}
// 	if err := json.Unmarshal(body, v); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (c *Clint) GetCountents(url *string) ([]git.TreeElement, error) {
// 	if c.GitRepoUrl == "" {
// 		return nil, errors.New("GitRepoUrl not set")
// 	}
// 	if url == nil {
// 		url = &c.GitRepoUrl
// 	}
// 	req, err := http.NewRequest("GET", *url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Accept", "application/vnd.github+json")
// 	if c.GitToken != "" {
// 		req.Header.Add("Authorization", "Bearer "+c.GitToken)
// 	}
// 	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
//
// 	var contents git.Tree
//
// 	err = c.sendRequest(req, &contents)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return contents.Tree, nil
// }
