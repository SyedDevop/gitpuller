package git_test

import (
	"testing"

	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/stretchr/testify/assert"
)

func TestLinkParser(t *testing.T) {
	out := []*gituser.Link{
		{
			Url: "https://api.github.com/user/89797705/repos?per_page=20&page=2",
			Rel: "next",
		},
		{
			Url: "https://api.github.com/user/89797705/repos?per_page=20&page=2",
			Rel: "last",
		},
	}

	rawLink := `<https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next", <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="last"`
	assert.Equal(t, out, gituser.ParseLinkHeader(rawLink))

	rawLink = ` <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next", <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="last" `
	assert.Equal(t, out, gituser.ParseLinkHeader(rawLink))

	rawLink = `  <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next"  , <https://api.github.com/user/89797705/repos?per_page=20&page=2>;  rel="last" `
	assert.Equal(t, out, gituser.ParseLinkHeader(rawLink))

	rawLink = ""
	assert.Nil(t, gituser.ParseLinkHeader(rawLink))

	rawLink = " "
	assert.Nil(t, gituser.ParseLinkHeader(rawLink))
}
