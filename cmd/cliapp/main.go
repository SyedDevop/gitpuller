package cliapp

import (
	"github.com/urfave/cli/v2"
)

func CliApp() *cli.App {
	app := &cli.App{
		Name:    "Git Puller",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "Syed uzair ahmed",
				Email: "syeds.devops007@gmail.com",
			},
		},
		UsageText: "git-puller [global options] command [command options] [arguments...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "branch",
				Usage:   "Specify the branch to pull from",
				Aliases: []string{"b"},
			},
		},
		// Action: func(c *cli.Context) error {
		// 	fmt.Println(c.Args().First())
		// 	return nil
		// },
	}
	return app
}

// func CliAppInit() {
// 	app := cliApp()
// 	if err := app.Run(os.Args); err != nil {
// 		log.Fatal(err)
// 	}
// }
