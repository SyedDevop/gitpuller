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
			FolderRepo:   make([]types.Repo, 0),
			RootPath:     baseFileName,
			CurPath:      baseFileName,
		}
		quitSelect := false

		t := tea.NewProgram(multiSelect.InitialModelMultiSelect(fetch, conTree, "Select File/Dir to download", &quitSelect))
		if _, err := t.Run(); err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		if fetch.Err != nil {
			log.Fatal(fetch.Err.Error())
		}
		if quitSelect || len(conTree.SelectedRepo) <= 0 {
			fmt.Println("\nNo option chosen 😊 Feel free to explore again!")
			os.Exit(0)
		}

		wg := sync.WaitGroup{}
		// st := tea.NewProgram(spinner.InitialModelNew("Processing... File to be downloaded..."))
		//
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	if _, err := st.Run(); err != nil {
		// 		log.Fatal(err)
		// 	}
		// }()
		// err = FetchAllFolders(conTree, fetch)
		// if err != nil {
		// 	if releaseErr := st.ReleaseTerminal(); releaseErr != nil {
		// 		log.Printf("Problem releasing terminal: %v", releaseErr)
		// 	}
		// 	return err
		// }
		// st.Quit()
		// err = st.ReleaseTerminal()
		// if err != nil {
		// 	if err != nil {
		// 		log.Printf("Could not release terminal: %v", err)
		// 		return err
		// 	}
		// }
		dt := tea.NewProgram(progress.InitialProgress(conTree.SelectedRepo))

		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dt.Run(); err != nil {
				log.Fatal(err)
			}
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

		fmt.Println("Execution Time: ", time.Since(start))
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func FetchAllFolders(conTree *multiSelect.ContentTree, fetch *multiSelect.Fetch) error {
	errChan := make(chan error, len(conTree.FolderRepo))
	var wg sync.WaitGroup

	for _, repo := range conTree.FolderRepo {
		wg.Add(1)
		go func(repo types.Repo) {
			defer wg.Done()
			allRepos, err := FetchRepoFiles(repo.URL, fetch)
			if err != nil {
				errChan <- err
				return
			}

			// Safely append to SelectedRepo
			conTree.Mu.Lock()
			conTree.SelectedRepo = append(conTree.SelectedRepo, allRepos...)
			conTree.Mu.Unlock()
		}(repo)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func FetchRepoFiles(url string, fetch *multiSelect.Fetch) ([]types.Repo, error) {
	var repos []types.Repo
	fetch.Clint.GitRepoUrl = url
	data, err := fetch.Clint.GetCountents(nil)
	if err != nil {
		return nil, err
	}
	rawData := util.GetRepoFromContent(*data)

	for _, item := range rawData {
		if item.Type == "dir" {
			// Recursively fetch contents from the directory, excluding the directory itself
			newData, err := FetchRepoFiles(item.URL, fetch)
			if err != nil {
				return nil, err
			}
			repos = append(repos, newData...)
		} else {
			// Append non-directory items to the list
			repos = append(repos, item)
		}
	}
	return repos, nil
}
