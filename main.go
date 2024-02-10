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
			fmt.Println("Please provide RepoName and UserName url example: gitpuller 'SyedDevop/gitpuller'")
			return nil
		}
		headderMes := fmt.Sprintf("Fetching your contents Form %s Repo", c.Args().Get(0))
		clint.GitRepoUrl = util.ParseContentsUrl(c.Args().Get(0))

		// Manager for Fetching State of git repo contents.
		fetch := &multiSelect.Fetch{
			Clint:     clint,
			FethDone:  false,
			FetchMess: headderMes,
		}
		conTree := &multiSelect.ContentTree{
			Tree:         make(map[string]*multiSelect.Node),
			SelectedRepo: make([]types.Repo, 0),
			RootPath:     "home",
			CurPath:      "home",
		}
		quitSelect := false

		t := tea.NewProgram(multiSelect.InitialModelMultiSelect(fetch, conTree, "Select File/Dir to download", &quitSelect))
		if _, err := t.Run(); err != nil {
			log.Fatal(err)
		}

		if fetch.Err != nil {
			log.Fatal(fetch.Err.Error())
		}
		if quitSelect {
			fmt.Println("\nNo option chosen 😊 Feel free to explore again!")
			os.Exit(0)
		}

		fmt.Println(conTree.Tree)
		dt := tea.NewProgram(progress.InitialProgress(conTree.SelectedRepo))

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dt.Run(); err != nil {
				log.Fatal(err)
			}
		}()

		for _, choice := range conTree.SelectedRepo {
			switch choice.Type {
			case "dir":
				fmt.Println("Directory Currently not supported")
			case "file":
				err := progress.DownloadFile(choice, "temp")
				if err != nil {
					releaseErr := dt.ReleaseTerminal()
					if releaseErr != nil {
						log.Fatalf("Problem releasing terminal: %v", releaseErr)
					}
				}

				dt.Send(progress.DownloadMes(choice.Name))
			}
		}

		dt.Quit()
		err := dt.ReleaseTerminal()
		if err != nil {
			log.Fatalf("Could not release terminal: %v", err)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
