package gituser

import (
	"fmt"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/assert"
	"github.com/SyedDevop/gitpuller/pkg/client"
)

type (
	Repos struct {
		NextLink    *string
		LastLink    *string
		PrevLink    *string
		FirstLink   *string
		Client      *client.Client
		CurrentPage int
		PageCount   int
		ItraterDone bool
	}

	Link struct {
		Url string
		Rel string
	}
	// ReposLinkIterator interface {
	// 	Next() ReposLink
	// }
	GitUser struct {
		Repos    *Repos
		Client   *client.Client
		UserName string
		UserUrl  string
	}
)

func NewGitUser() *GitUser {
	c := client.NewClint()

	c.AddHeader("Accept", "application/vnd.github+json")
	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	gitToken := util.GetGitToken()
	if gitToken != "" {
		c.AddBareAuth(gitToken)
	}

	repos := &Repos{
		CurrentPage: 1,
		PageCount:   0,
		Client:      c,
		ItraterDone: false,
	}

	return &GitUser{
		Repos:    repos,
		Client:   c,
		UserName: "",
		UserUrl:  "",
	}
}

func (g *GitUser) SetUserUrl(url string)   { g.UserUrl = url }
func (g *GitUser) SetUserName(name string) { g.UserName = name }
func (g *GitUser) ProjectName() string     { return g.UserName }
func (g *GitUser) Name() string            { return g.UserName }
func (g *GitUser) Description() string {
	return "This is the Repositories of " + g.UserUrl + ""
}

func (g *GitUser) GetUsersRepos(url string) ([]UserRepos, error) {
	res, err := g.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var userRepos []UserRepos
	err = client.UnmarshalJSON(res, &userRepos)
	if err != nil {
		return nil, err
	}
	return userRepos, nil
}

func (r *Repos) SetNextLink(url string) {
	r.NextLink = &url
}

func (r *Repos) Reset() {
	r.ItraterDone = false
	r.NextLink = r.FirstLink
	r.FirstLink = nil
	r.PrevLink = nil
	r.LastLink = nil
}

func (r *Repos) String() string {
	n, f, p, l := "Empty no Url", "Empty no Url", "Empty no Url", "Empty no Url"

	if r.FirstLink != nil {
		f = *r.FirstLink
	}
	if r.NextLink != nil {
		n = *r.NextLink
	}
	if r.PrevLink != nil {
		p = *r.PrevLink
	}
	if r.LastLink != nil {
		l = *r.LastLink
	}

	return fmt.Sprintf(`
Next Link: %s,
Prev Link: %s,
First Link: %s,
Last Link: %s`, n, p, f, l)
}

func (r *Repos) Next() ([]UserRepos, error) {
	assert.Assert(r.NextLink != nil, "The next url link for UserRepos is nil")
	if r.ItraterDone {
		return nil, nil
	}
	res, err := r.Client.Get(*r.NextLink)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body := make([]byte, 1024)
		res.Body.Read(body)
		return nil, fmt.Errorf("failed to get repos %s", string(body))
	}

	links := ParseLinkHeader(res.Header.Get("Link"))
	linksLen := len(links)

	for _, link := range links {
		switch link.Rel {
		case "next":
			r.NextLink = &link.Url
		case "last":
			r.LastLink = &link.Url
		case "first":
			if linksLen == 2 {
				r.ItraterDone = true
			}
			r.FirstLink = &link.Url
		case "prev":
			if linksLen == 2 {
				r.ItraterDone = true
			}
			r.PrevLink = &link.Url
		}
	}

	var userRepos []UserRepos
	err = client.UnmarshalJSON(res, &userRepos)
	if err != nil {
		return nil, err
	}
	return userRepos, nil
}
