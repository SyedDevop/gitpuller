package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SyedDevop/gitpuller/api"
	"github.com/SyedDevop/gitpuller/cliapp"
	"github.com/SyedDevop/gitpuller/ui/multiSelect"
	"github.com/SyedDevop/gitpuller/ui/progress"
	"github.com/SyedDevop/gitpuller/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	types "github.com/SyedDevop/gitpuller/mytypes"
)

func main() {
	app := cliapp.CliApp()
	clint := api.NewClint()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access environment variables
	gitToken := os.Getenv("GIT_TOKEN")
	clint.GitToken = gitToken

	app.Action = func(c *cli.Context) error {
		if c.NArg() <= 0 {
			fmt.Println("Please provide RepoName and UserName url example: gitpuller 'SyedDevop/gitpuller'")
			return nil
		}

		contentUrl := c.Args().Get(0)
		headderMes := fmt.Sprintf("Fetching your contents Form %s Repo", contentUrl)
		clint.GitRepoUrl = util.ParseContentsUrl(contentUrl)

		baseFileName := strings.Split(contentUrl, "/")[1]
		// Manager for Fetching State of git repo contents.
		fetch := &multiSelect.Fetch{
			Clint:     clint,
			FethDone:  false,
			FetchMess: headderMes,
		}
		conTree := &multiSelect.ContentTree{
			Tree:         make(map[string]*multiSelect.Node),
			SelectedRepo: make([]types.Repo, 0),
			RootPath:     baseFileName,
			CurPath:      baseFileName,
		}
		quitSelect := false

		t := tea.NewProgram(multiSelect.InitialModelMultiSelect(fetch, conTree, "Select File/Dir to download", &quitSelect))
		if _, err := t.Run(); err != nil {
			log.Fatal(err)
		}

		if fetch.Err != nil {
			log.Fatal(fetch.Err.Error())
		}
		if quitSelect || len(conTree.SelectedRepo) <= 0 {
			fmt.Println("\nNo option chosen 😊 Feel free to explore again!")
			os.Exit(0)
		}
		dt := tea.NewProgram(progress.InitialProgress(conTree.SelectedRepo))

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dt.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		start := time.Now()
		defer func() {
			fmt.Println("Execution Time: ", time.Since(start))
		}()
		wg.Add(len(conTree.SelectedRepo))
		for _, choice := range conTree.SelectedRepo {
			go func(repo types.Repo) {
				defer wg.Done()
				err := progress.DownloadFile(repo, baseFileName)
				if err != nil {
					releaseErr := dt.ReleaseTerminal()
					if releaseErr != nil {
						log.Fatalf("Problem releasing terminal: %v", releaseErr)
					}
					log.Fatalf("Error while downloading %v", err)
				}

				dt.Send(progress.DownloadMes(repo.Name))
			}(choice)
		}
		wg.Wait()
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
