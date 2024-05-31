package gituser

import (
	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
)

type (
	ReposLink struct {
		NextLink, LastLink, PrevLink, FirstLink *string
		CurrentPage                             int
		PageCount                               int
	}
	Link struct {
		Url string
		Rel string
	}
	// ReposLinkIterator interface {
	// 	Next() ReposLink
	// }
	GitUser struct {
		ReposLink *ReposLink
		Client    *client.Client
		Name      string
	}
)

func (u UserRepos) Title() string { return u.Name }
func (u UserRepos) Description() string {
	if u.Descript == nil {
		return "No Description"
	}
	return *u.Descript
}
func (u UserRepos) FilterValue() string { return u.Name }

func NewGitUser(name string) *GitUser {
	c := client.NewClint()

	c.AddHeader("Accept", "application/vnd.github+json")
	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	gitToken := util.GetGitToken()
	if gitToken != "" {
		c.AddBareAuth(gitToken)
	}

	repoLinlk := &ReposLink{
		CurrentPage: 1,
		PageCount:   0,
	}

	return &GitUser{
		ReposLink: repoLinlk,
		Client:    c,
		Name:      name,
	}
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

func (r *ReposLink) Next() *ReposLink {
	panic("Implement me the next iterator for ReposLink")
}
