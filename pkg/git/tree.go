package git

import (
	"io/fs"
)

type BadReq struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

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
	Type FileType `json:"type"`
	SHA  string   `json:"sha"`
	Mode FileMode `json:"mode,string"`
}

// Type.go

type FileType string

const (
	Blob     FileType = "blob"
	TypeTree FileType = "tree"
)

func (t *TreeElement) IsTree() bool {
	return t.Type == TypeTree
}

func (t *TreeElement) TreeType() string {
	switch t.Type {
	case Blob:
		return "file"
	case TypeTree:
		return "dir"
	default:
		return ""
	}
}

type Blobl struct {
	SHA      string `json:"sha"`
	NodeID   string `json:"node_id"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	Size     int64  `json:"size"`
}

type FileMode int

const (
	FileModeEmpty      FileMode = 0
	FileModeTree       FileMode = 40000
	FileModeBlob       FileMode = 100644
	FileModeExec       FileMode = 100755
	FileModeSymlink    FileMode = 120000
	FileModeCommit     FileMode = 160000
	FileModeDeprecated FileMode = 100664
)

func ToOSFileMode(m FileMode) fs.FileMode {
	switch m {
	case FileModeTree:
		return fs.ModePerm | fs.ModeDir
	case FileModeCommit:
		return fs.ModePerm | fs.ModeDir
	case FileModeBlob:
		return fs.FileMode(0644)
	// Deprecated is no longer awed: treated as a Regular instead
	case FileModeDeprecated:
		return fs.FileMode(0644)
	case FileModeExec:
		return fs.FileMode(0755)
	case FileModeSymlink:
		return fs.ModePerm | fs.ModeSymlink
	}

	return fs.FileMode(0)
}
