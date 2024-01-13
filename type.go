package main

type Content struct {
	DownloadURL *string `json:"download_url,omitempty"`
	Links       Links   `json:"_links"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	SHA         string  `json:"sha"`
	URL         string  `json:"url"`
	HTMLURL     string  `json:"html_url"`
	GitURL      string  `json:"git_url"`
	Type        string  `json:"type"`
	Content     string  `json:"content"`
	Encoding    string  `json:"encoding"`
	Size        int64   `json:"size"`
}

type Links struct {
	Self string `json:"self"`
	Git  string `json:"git"`
	HTML string `json:"html"`
}
