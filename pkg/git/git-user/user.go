package gituser

import (
	"fmt"
	"time"
)

type UserRepos struct {
	Descript        *string   `json:"description"`
	License         *License  `json:"license"`
	Language        *string   `json:"language"`
	Homepage        *string   `json:"homepage"`
	SSHURL          string    `json:"ssh_url"`
	CommitsURL      string    `json:"commits_url"`
	HTMLURL         string    `json:"html_url"`
	DefaultBranch   string    `json:"default_branch"`
	Visibility      string    `json:"visibility"`
	URL             string    `json:"url"`
	ForksURL        string    `json:"forks_url"`
	BranchesURL     string    `json:"branches_url"`
	TagsURL         string    `json:"tags_url"`
	BlobsURL        string    `json:"blobs_url"`
	GitTagsURL      string    `json:"git_tags_url"`
	GitRefsURL      string    `json:"git_refs_url"`
	TreesURL        string    `json:"trees_url"`
	LanguagesURL    string    `json:"languages_url"`
	ContributorsURL string    `json:"contributors_url"`
	NodeID          string    `json:"node_id"`
	ContentsURL     string    `json:"contents_url"`
	DownloadsURL    string    `json:"downloads_url"`
	LabelsURL       string    `json:"labels_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	GitURL          string    `json:"git_url"`
	Name            string    `json:"name"`
	CloneURL        string    `json:"clone_url"`
	FullName        string    `json:"full_name"`
	Owner           Owner     `json:"owner"`
	Size            int64     `json:"size"`
	ID              int64     `json:"id"`
	AllowForking    bool      `json:"allow_forking"`
	Fork            bool      `json:"fork"`
	Private         bool      `json:"private"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type Owner struct {
	Login      string `json:"login"`
	NodeID     string `json:"node_id"`
	AvatarURL  string `json:"avatar_url"`
	GravatarID string `json:"gravatar_id"`
	URL        string `json:"url"`
	HTMLURL    string `json:"html_url"`
	ReposURL   string `json:"repos_url"`
	Type       string `json:"type"`
	ID         int64  `json:"id"`
	SiteAdmin  bool   `json:"site_admin"`
}

func (u UserRepos) Command() string { return "git clone " + u.CloneURL }
func (u UserRepos) IsPrivate() bool { return u.Private }
func (u UserRepos) Title() string   { return u.Name }
func (u UserRepos) Description() string {
	if u.Descript == nil {
		return "No Description"
	}
	return *u.Descript
}
func (u UserRepos) FilterValue() string { return u.Name }

func (u UserRepos) String() string {
	dis := ""
	if u.Descript != nil {
		dis = *u.Descript
	}
	return fmt.Sprintf(`
Repo Name: %s
Full Repo Name: %s
Owner: %s
Discretion %v
    `, u.Name, u.FullName, u.Owner.Login, dis)
}
