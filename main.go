package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SyedDevop/gitpuller/cliapp"
	"github.com/SyedDevop/gitpuller/ui/loader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

// _ "github.com/SyedDevop/gitpuller/cliapp"

func main() {
	app := cliapp.CliApp()
	clint := NewClint()
	app.Action = func(c *cli.Context) error {
		if c.NArg() <= 0 {
			fmt.Println("Please provide repo url")
			return nil
		}
		loadSatte := loader.InitialModelSpiner("Fetching Git Repo")
		p := tea.NewProgram(loadSatte)
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Welcome to Git Puller")
		path := c.Args().Get(0)
		contents, err := clint.getCountents(path)
		if err != nil {
			return err
		}
		fmt.Println(contents)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
