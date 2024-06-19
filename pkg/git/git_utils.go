package git

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/SyedDevop/gitpuller/pkg/assert"
)

const (
	LOCALHOST = "http://localhost:4069"
	HOST      = "https://api.github.com"
)

func GetDomain() string {
	dev := os.Getenv("DEV")
	if dev == "LOCAL" {
		return LOCALHOST
	}
	return HOST
}

func CheckDomain(url string) string {
	if strings.Contains(url, HOST) {
		return strings.Replace(url, HOST, LOCALHOST, 1)
	}
	return url
}

// UserReposURL generates a GitHub API URL to fetch repositories for a given user.
// returns the formatted URL as a string.
//
// Parameters:
//   - name: The GitHub username as a string.
//
// Returns:
//   - A formatted URL string to access the user's repositories on GitHub.
func UserReposURL(name string) string {
	return fmt.Sprintf("%s/users/%s/repos", GetDomain(), strings.TrimSpace(name))
}

// AuthReposURL generates a GitHub API URL to fetch repositories for Authenticated user.
// returns the formatted URL as a string.
//
// Returns:
//   - A formatted URL string to access the Authenticated user's repositories on GitHub.
func AuthReposURL() string {
	return fmt.Sprintf("%s/user/repos", GetDomain())
}

// Function to check if a string contains at least one letter
func hasLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// Function to check if both sides of '/' contain words
func hasWords(s string) bool {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return false
	}
	return hasLetters(parts[0]) && hasLetters(parts[1])
}

func RepoUrl(repoCredential string, recursice bool) string {
	assert.Assert(hasWords(repoCredential), "RepoUrl Credential are not in proper formate exc-formate:'SyedDevop/gitpuller' got:", repoCredential)
	rawUrl := strings.Builder{}
	rawUrl.WriteString(GetDomain())
	rawUrl.WriteString(fmt.Sprintf("/repos/%s/git/trees/main", repoCredential))
	if recursice {
		rawUrl.WriteString("?recursice=1")
	}
	return rawUrl.String()
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
