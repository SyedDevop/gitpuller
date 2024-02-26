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

func DownloadFile(content *api.Repo, rootPath string) error {
	// Get the download URL
	downloadURL := content.DownloadURL
	if downloadURL == nil {
		// log.Fatal("The Download URL is not available")
		return errors.New("download URL not available")
	}

	// Get the data
	resp, err := http.Get(*downloadURL)
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
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
