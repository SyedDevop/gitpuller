package git

import (
	"fmt"
	"os"
	"strings"
)

// UserReposURL generates a GitHub API URL to fetch repositories for a given user.
// returns the formatted URL as a string.
//
// Parameters:
//   - name: The GitHub username as a string.
//
// Returns:
//   - A formatted URL string to access the user's repositories on GitHub.
func UserReposURL(name string) string {
	dev := os.Getenv("DEV")
	host := "https://api.github.com"
	if dev == "LOCAL" {
		host = "http://localhost:4069"
	}
	return fmt.Sprintf("%s/users/%s/repos", host, strings.TrimSpace(name))
}

// AuthReposURL generates a GitHub API URL to fetch repositories for Authenticated user.
// returns the formatted URL as a string.
//
// Returns:
//   - A formatted URL string to access the Authenticated user's repositories on GitHub.
func AuthReposURL() string {
	dev := os.Getenv("DEV")
	host := "https://api.github.com"
	if dev == "LOCAL" {
		host = "http://localhost:4069"
	}
	return fmt.Sprintf("%s/user/repos", host)
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

	return fmt.Sprintf("%s?per_page=%d&page=%d", url, *per, *pages)
}
