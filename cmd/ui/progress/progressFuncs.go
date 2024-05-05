package progress

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SyedDevop/gitpuller/cmd/api"
	"github.com/SyedDevop/gitpuller/cmd/util"
)

// FIX: file permissions. Use the provided permission from the git it self.
// TODO : Choose whether get the file as raw or blob.
func DownloadFile(content *api.TreeElement, rootPath string) error {
	// Get the download URL
	downloadURL := content.URL
	if downloadURL == nil {
		// log.Fatal("The Download URL is not available")
		return errors.New("download URL not available")
	}
	req, err := http.NewRequest("GET", *downloadURL, nil)

	gitToken := os.Getenv("GIT_TOKEN")
	if gitToken != "" {
		req.Header.Add("Authorization", "Bearer "+gitToken)
	}
	req.Header.Add("Accept", "application/vnd.github.raw+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(rootPath, content.Path)

	isRoot, dirPath := util.GetParentPath(filePath)
	if !isRoot {
		if err := util.CreateDir(dirPath); err != nil {
			return err
		}
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = out.Chmod(api.ToOSFileMode(content.Mode))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
