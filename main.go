package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/SyedDevop/gitpuller/api"
	"github.com/SyedDevop/gitpuller/cliapp"
	"github.com/SyedDevop/gitpuller/ui/multiSelect"
	"github.com/SyedDevop/gitpuller/ui/progress"
	"github.com/SyedDevop/gitpuller/ui/spinner"
	"github.com/SyedDevop/gitpuller/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	types "github.com/SyedDevop/gitpuller/mytypes"
)

func main() {
	app := cliapp.CliApp()
	clint := api.NewClint()
	app.Action = func(c *cli.Context) error {
		if c.NArg() <= 0 {
			fmt.Println("Please provide types.Repo url")
			return nil
		}
		headderMes := fmt.Sprintf("Fetching your contents Form %s Repo", c.Args().Get(0))
		spinner := tea.NewProgram(spinner.InitialModelNew(headderMes))

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
		contents, err := clint.GetCountents(path)
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

		repos := util.GetRepoFromContent(*contents)

		sel := &multiSelect.Selection{
			Choices: make([]types.Repo, 0),
		}
		quitSelect := false

		t := tea.NewProgram(multiSelect.InitialModelMultiSelect(repos, sel, "Select File/Dir to download", &quitSelect))
		if _, err := t.Run(); err != nil {
			log.Fatal(err)
		}

		if quitSelect {
			fmt.Println("\nNo option chosen 😊 Feel free to explore again!")
			os.Exit(0)
		}

		dt := tea.NewProgram(progress.InitialProgress(sel.Choices))

		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dt.Run(); err != nil {
				log.Fatal(err)
			}
		}()

		for _, choice := range sel.Choices {
			switch choice.Type {
			case "dir":
				fmt.Println("Directory Currently not supported")
			case "file":
				// fmt.Println(choice.Name)
				if err = progress.DownloadFile(choice, "temp"); err != nil {
					releaseErr := dt.ReleaseTerminal()
					if releaseErr != nil {
						log.Fatalf("Problem releasing terminal: %v", releaseErr)
					}

				}
				dt.Send(progress.DownloadMes(choice.Name))
			}
		}

		dt.Quit()
		err = dt.ReleaseTerminal()
		if err != nil {
			log.Fatalf("Could not release terminal: %v", err)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
