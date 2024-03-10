package api

// type Content struct {
// 	DownloadURL *string `json:"download_url,omitempty"`
// 	Links       Links   `json:"_links"`
// 	Name        string  `json:"name"`
// 	Path        string  `json:"path"`
// 	SHA         string  `json:"sha"`
// 	URL         string  `json:"url"`
// 	HTMLURL     string  `json:"html_url"`
// 	GitURL      string  `json:"git_url"`
// 	Type        string  `json:"type"`
// 	Content     string  `json:"content"`
// 	Encoding    string  `json:"encoding"`
// 	Size        uint64  `json:"size"`
// }

type BadReq struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

// type Links struct {
// 	Self string `json:"self"`
// 	Git  string `json:"git"`
// 	HTML string `json:"html"`
// }

// Tree.go

type Tree struct {
	SHA       string        `json:"sha"`
	URL       string        `json:"url"`
	Tree      []TreeElement `json:"tree"`
	Truncated bool          `json:"truncated"`
}

// TreeElement.go

type TreeElement struct {
	Size *int64   `json:"size,omitempty"`
	URL  *string  `json:"url"`
	Path string   `json:"path"`
	Mode string   `json:"mode"`
	Type FileType `json:"type"`
	SHA  string   `json:"sha"`
}

// Type.go

type FileType string

const (
	Blob     FileType = "blob"
	TypeTree FileType = "tree"
)

type Blobl struct {
	SHA      string `json:"sha"`
	NodeID   string `json:"node_id"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	Size     int64  `json:"size"`
}

// type Repo struct {
// 	DownloadURL *string
// 	Type        string
// 	Name        string
// 	Path        string
// 	URL         string
// 	Size        uint64
// }
