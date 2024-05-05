/*
Copyright © 2024 Syed Uzair Ahmed <syeds.devops007@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SyedDevop/gitpuller/cmd/api"
	"github.com/SyedDevop/gitpuller/cmd/ui/multiSelect"
	"github.com/SyedDevop/gitpuller/cmd/ui/progress"
	"github.com/SyedDevop/gitpuller/cmd/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var parentFlag bool

var getCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Get the file/folder from remote Git repository",
	Long: `Get the file/folder from remote Git repository

Example: gitpuller get SyedDevop/gitpuller
  `,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			cobra.CheckErr("Please provide RepoName and UserName url example: gitpuller get 'SyedDevop/gitpuller'")
		}
		clint := api.NewClint()
		contentUrl := args[0]
		headderMes := fmt.Sprintf("Fetching your contents Form %s Repo", contentUrl)
		clint.GitRepoUrl = util.ParseContentsUrl(contentUrl, "main")

		rootPath := ""
		if len(args) == 2 {
			rootPath = args[1]
		}

		urlFilePath := strings.Split(contentUrl, "/")[1]

		// Manager for Fetching State of git repo contents.
		fetch := &multiSelect.Fetch{
			Clint:     clint,
			FethDone:  false,
			FetchMess: headderMes,
		}
		conTree := &multiSelect.ContentTree{
			Tree:         make(map[string]*multiSelect.Node),
			SelectedRepo: make(map[string][]api.TreeElement),
			FolderRepo:   make([]api.TreeElement, 0),
			RootPath:     urlFilePath,
			CurPath:      urlFilePath,
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

		wg := sync.WaitGroup{}

		dt := tea.NewProgram(progress.InitialProgress(conTree.SelectedRepoLen()))
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := dt.Run(); err != nil {
				log.Fatal(err)
			}
		}()

		wg.Add(conTree.SelectedRepoLen())
		for filePath, choice := range conTree.SelectedRepo {
			go func(repos []api.TreeElement, repoPath string) {
				for _, repo := range repos {
					go func(repo *api.TreeElement) {
						defer wg.Done()
						path := filepath.Join(rootPath, repoPath)
						// fmt.Printf("downloading:: to %s :: Data %+v \n", path, repo)
						err := progress.DownloadFile(repo, path)
						if err != nil {
							releaseErr := dt.ReleaseTerminal()
							if releaseErr != nil {
								log.Fatalf("Problem releasing terminal: %v", releaseErr)
							}
							log.Fatalf("Error while downloading %v", err)
						}
						dt.Send(progress.DownloadMes(repo.Path))
					}(&repo)
				}
			}(choice, filePath)
		}

		wg.Wait()
		dt.Quit()
		err := dt.ReleaseTerminal()
		if err != nil {
			log.Fatalf("Could not release terminal: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.Flags().BoolVarP(&parentFlag, "parent", "p", false, "Keep the repository's parent directory")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
