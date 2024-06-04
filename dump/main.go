package main

import (
	"flag"
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

func run(path, fileName string) error {
	c := client.NewClint()
	c.AddHeader("Accept", "application/vnd.github+json")
	c.AddHeader("X-GitHub-Api-Version", "2022-11-28")
	gitToken := util.GetGitToken()
	if gitToken != "" {
		c.AddBareAuth(gitToken)
	}
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
	fName := filepath.Join(path, fileName)
	file, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Print("Done downloading file", "path", path, "fileName", fileName, "duration", time.Since(start))
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

var (
	path = flag.String("path", ".", "the path to dump the file")
	file = flag.String("file", "file.json", "File name to dump")
)

func main() {
	flag.Parse()
	if err := run(*path, *file); err != nil {
		log.Fatal("Error from Dump main", "err", err)
	}
}
