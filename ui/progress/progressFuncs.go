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

func DownloadFile(content types.Repo, dest string) error {
	// fmt.Println("Downloading:", content.Name)

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

	if dest != "" {
		util.CreateDir(dest)
	}

	// Create the file
	out, err := os.Create(filepath.Join(dest, content.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Add delay for testing
	// time.Sleep(3 * time.Second)

	return nil
}
