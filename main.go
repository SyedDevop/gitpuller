package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/SyedDevop/gitpuller/cliapp"
	"github.com/SyedDevop/gitpuller/ui/multiSelect"
	"github.com/SyedDevop/gitpuller/ui/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	. "github.com/SyedDevop/gitpuller/mytypes"
)

func main() {
	app := cliapp.CliApp()
	clint := NewClint()
	app.Action = func(c *cli.Context) error {
		if c.NArg() <= 0 {
			fmt.Println("Please provide repo url")
			return nil
		}
		spinner := tea.NewProgram(spinner.InitialModelNew("Fetching your repo"))

		// add synchronization to wait for spinner to finish
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := spinner.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()

		path := c.Args().Get(0)
		contents, err := clint.getCountents(path)
		if err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				log.Printf("Problem releasing terminal: %v", releaseErr)
			}
			return err
		}
		err = spinner.ReleaseTerminal()
		if err != nil {
			log.Printf("Could not release terminal: %v", err)
			return err
		}

		repos := getRepoFromContent(*contents)

		sel := &multiSelect.Selection{
			Choices: make([]Repo, 0),
		}

		t := tea.NewProgram(multiSelect.InitialModelMultiSelect(repos, sel, "Select"))
		if _, err := t.Run(); err != nil {
			log.Fatal(err)
		}
		for _, choice := range sel.Choices {
			switch choice.Type {
			case "dir":
				fmt.Println("Directory Currently not supported")
			case "file":
				downloadFile(choice, "")
			}
		}

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func downloadFile(content Repo, dest string) {
	fmt.Println("Downloading:", content.Name)

	// Get the download URL
	downloadURL := content.DownloadURL
	if downloadURL == nil {
		log.Fatal("The Download URL is not available")
	}

	// Get the data
	resp, err := http.Get(*downloadURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath.Join(dest, content.Name))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
