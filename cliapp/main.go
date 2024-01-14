package cliapp

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func cliApp() *cli.App {
	app := &cli.App{
		Name:    "Git Puller",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "Syed uzair ahmed",
				Email: "syeds.devops007@gmail.com",
			},
		},
		HelpName:  "contrive",
		Usage:     "demonstrate available API",
		UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
	}
	return app
}

func CliAppInit() {
	app := cliApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
