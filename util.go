package main

import "fmt"

func parseContentsUrl(path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents", path)
}
