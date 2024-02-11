package progress

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	types "github.com/SyedDevop/gitpuller/mytypes"
	"github.com/SyedDevop/gitpuller/util"
)

func DownloadFile(content types.Repo, rootPath string) error {
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

	filePath := filepath.Join(rootPath, content.Name)
	if content.Name != content.Path {
		filePath = filepath.Join(rootPath, content.Path)
	}

	_, dirPath := util.GetParentPath(filePath)
	if err := util.CreateDir(dirPath); err != nil {
		return err
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
