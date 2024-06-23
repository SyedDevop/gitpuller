package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func getGitClient() *client.Client {
	c := client.NewClint()
	c.AddHeader("Accept", "application/vnd.github+json")
	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")
	gitToken := util.GetGitToken()
	if gitToken != "" {
		c.AddBareAuth(gitToken)
	}
	return c
}

func run(_, user string) {
	c := getGitClient()

	start := time.Now()
	log.Info("Start downloading Repos")
	repos, err := userRepos(c)
	if err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("%s.json", user)
	fileLocation := filepath.Join(basePath, fileName)
	log.Info("Done downloading Repos", "Path", fileLocation, "duration", time.Since(start))

	start = time.Now()
	log.Info("Start Downloading repo tree")
	_, errs := getGitFile(c, repos)
	if len(errs) != 0 {
		for _, err := range errs {
			log.Error("run#getGitFile", err)
		}
		os.Exit(1)
	}
	repoPath := filepath.Join(basePath, "repo")
	log.Info("Done downloading Repo Tree", "Path", repoPath, "duration", time.Since(start))
}

var (
	path     = flag.String("path", "", "(Required)::The path to dump the file in.")
	user     = flag.String("user", "", "(Required)::Github user name")
	h        = flag.Bool("h", false, "Show help")
	basePath = ""
)

func init() {
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	if flag.NFlag() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	basePath = filepath.Join(*path, *user)
}

func main() {
	run(*path, *user)
}
