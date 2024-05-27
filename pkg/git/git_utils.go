package git

import (
	"strings"

	"github.com/SyedDevop/gitpuller/pkg/assert"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
)

// Response:  <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next", <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="last"
func ParseLinkHeader(rawLink string) []*gituser.Link {
	if len(rawLink) < 2 {
		return nil
	}

	links := strings.Split(rawLink, ",")
	assert.Assert(len(links) == 2, "GitUser#ParseLinkHeader expected the link to have two links only go(", len(links), ")\n")
	linkList := make([]*gituser.Link, 2)
	for i, link := range links {
		data := strings.Split(link, ">;")
		assert.Assert(len(data) == 2, "GitUser#ParseLinkHeader::Url and Rel to be got(", data[0], data[1], ")\n")
		url := strings.TrimSpace(data[0])
		url = url[1:]
		rel := strings.TrimSpace(data[1])
		rel = rel[5 : len(rel)-1]
		linkList[i] = &gituser.Link{
			Url: url,
			Rel: rel,
		}
	}

	return linkList
}
