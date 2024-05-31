package gituser_test

import (
	"fmt"
	"testing"

	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/stretchr/testify/require"
)

func TestGetUserRepos(t *testing.T) {
	c := client.NewClint()

	c.AddHeader("Accept", "application/vnd.github+json")
	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")

	per := 5
	page := 1
	res, err := c.Get(git.AddPaginationParams(git.GenerateReposURL("SyedDevop"), &per, &page))
	require.NoError(t, err)
	defer res.Body.Close()

	var userRepos []gituser.UserRepos
	err = client.UnmarshalJSON(res, &userRepos)

	fmt.Println(userRepos)
	require.NoError(t, err)
}
