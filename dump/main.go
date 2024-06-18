package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/client"
	"github.com/SyedDevop/gitpuller/pkg/git"
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

func run(path, fileName, getData string) error {
	fileJson := fmt.Sprintf("%s.json", fileName)
	fName := filepath.Join(path, fileJson)
	if getData == "repo" {
		start := time.Now()
		log.Info("Start Downloading repo tree from", "file", fName)

		if errs := getGitFile(fName, "repo"); len(errs) != 0 {
			for _, err := range errs {
				log.Error("run#getGitFile", err)
			}
			os.Exit(1)
		}

		log.Print("Done Downloading repo tree from", "file", fName, "duration", time.Since(start))
		os.Exit(0)
	}

	c := getGitClient()
	per := 100
	pages := 1

	start := time.Now()
	log.Info("Start downloading file", "path", path, "fileName", fileName)

	res, err := c.Get(git.AddPaginationParams(git.AuthReposURL(), &per, &pages))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Info("Creating the file", "path", path, "fileName", fileName)
	// fileJson := fmt.Sprintf("%s.json", fileName)
	// fName := filepath.Join(path, fileJson)
	file, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Print("Done downloading file", "path", path, "fileName", fileJson, "duration", time.Since(start))
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

var (
	path = flag.String("path", "", "the path to dump the file")
	file = flag.String("file", "file.json", "File name to dump")
	data = flag.String("data", "user", "Type of data to fetch userRepos('user')/repo")
	h    = flag.Bool("h", false, "Show help")
)

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}
	if path == nil || *path == "" {
		fmt.Println("Path is required")
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*path, *file, *data); err != nil {
		log.Fatal("Error from Dump main", "err", err)
	}
}
