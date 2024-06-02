package gituser_test

import (
	"testing"

	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/stretchr/testify/require"
)

// func TestGetUserRepos(t *testing.T) {
// 	c := client.NewClint()
//
// 	c.AddHeader("Accept", "application/vnd.github+json")
// 	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")
//
// 	per := 5
// 	page := 1
// 	res, err := c.Get(git.AddPaginationParams(git.GenerateReposURL("SyedDevop"), &per, &page))
// 	require.NoError(t, err)
// 	defer res.Body.Close()
//
// 	var userRepos []gituser.UserRepos
// 	err = client.UnmarshalJSON(res, &userRepos)
//
// 	fmt.Println(userRepos)
// 	require.NoError(t, err)
// }

func TestGetUserNext(t *testing.T) {
	gituser := gituser.NewGitUser("SyedDevop")

	per := 5
	page := 1
	gituser.ReposLink.SetNextLink(git.AddPaginationParams(git.UserReposURL("SyedDevop"), &per, &page))

	count := 0
OUTER:
	for {
		val, err := gituser.ReposLink.Next()
		require.NoError(t, err)
		if val == nil {
			break OUTER
		}
		count++
	}
	require.EqualValues(t, 5, count)
}
