package gituser

import (
	"strings"

	"github.com/SyedDevop/gitpuller/pkg/assert"
)

func ParseLinkHeader(rawLink string) []*Link {
	// Response:  <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next", <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="last"
	if len(rawLink) < 2 {
		return nil
	}

	links := strings.Split(rawLink, ",")
	linkLenAssert := len(links) >= 2 && len(links) <= 4
	assert.Assert(linkLenAssert, "GitUser#ParseLinkHeader expected the link to have two or four links only got::", len(links), "\n\nRaw Link::", rawLink, "\n\nSplit Link::", links, "\n")
	linkList := make([]*Link, len(links))
	for i, link := range links {
		data := strings.Split(link, ">;")
		assert.Assert(len(data) == 2, "GitUser#ParseLinkHeader::Url and Rel to be got(", data[0], data[1], ")\n")
		url := strings.TrimSpace(data[0])
		url = url[1:]
		rel := strings.TrimSpace(data[1])
		rel = rel[5 : len(rel)-1]
		linkList[i] = &Link{
			Url: url,
			Rel: rel,
		}
	}

	return linkList
}
