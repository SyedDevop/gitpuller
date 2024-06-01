package git

import (
	"fmt"
	"strings"
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
