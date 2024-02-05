package mytypes

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
	Size        uint64  `json:"size"`
}

type BadReq struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

type Links struct {
	Self string `json:"self"`
	Git  string `json:"git"`
	HTML string `json:"html"`
}

type Repo struct {
	DownloadURL *string
	Type        string
	Name        string
	Path        string
	Size        uint64
}
