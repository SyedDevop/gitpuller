package git

import (
	"fmt"
	"strings"

	"github.com/SyedDevop/gitpuller/pkg/assert"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
)

// GenerateReposURL generates a GitHub API URL to fetch repositories for a given user.
// It takes a user's name as input, trims any leading or trailing whitespace, and
// returns the formatted URL as a string.
//
// Parameters:
//   - name: The GitHub username as a string.
//
// Returns:
//   - A formatted URL string to access the user's repositories on GitHub.
func GenerateReposURL(name string) string {
	return fmt.Sprintf("https://api.github.com/users/%s/repos", strings.TrimSpace(name))
}

// AddPaginationParams adds pagination parameters to a given URL.
//
// It takes a base URL and pointers to the number of items per page and the number of pages.
// If the provided per or pages parameters are nil, default values of 50 and 1 are used respectively.
//
// Parameters:
//   - url: The base URL as a string.
//   - per: A pointer to an integer specifying the number of items per page.
//   - pages: A pointer to an integer specifying the number of pages.
//
// Returns:
//   - A formatted URL string with the pagination parameters appended.
func AddPaginationParams(url string, per, pages *int) string {
	if per == nil {
		defaultPer := 50
		per = &defaultPer
	}

	if pages == nil {
		defaultPages := 1
		pages = &defaultPages
	}

	return fmt.Sprintf("%s?per_page=%d&pages=%d", url, *per, *pages)
}

func ParseLinkHeader(rawLink string) []*gituser.Link {
	// Response:  <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="next", <https://api.github.com/user/89797705/repos?per_page=20&page=2>; rel="last"
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
