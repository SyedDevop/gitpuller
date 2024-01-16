package main

import (
	"fmt"
	"log"

	"github.com/SyedDevop/gitpuller/multiSelect"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// app := cliapp.CliApp()
	// clint := NewClint()
	// app.Action = func(c *cli.Context) error {
	// 	if c.NArg() <= 0 {
	// 		fmt.Println("Please provide repo url")
	// 		return nil
	// 	}
	// 	spinner := tea.NewProgram(spinner.InitialModelNew("Fetching your repo"))
	//
	// 	// add synchronization to wait for spinner to finish
	// 	wg := sync.WaitGroup{}
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		if _, err := spinner.Run(); err != nil {
	// 			fmt.Println(err)
	// 			os.Exit(1)
	// 		}
	// 	}()
	//
	// 	path := c.Args().Get(0)
	// 	contents, err := clint.getCountents(path)
	// 	if err != nil {
	// 		if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
	// 			log.Printf("Problem releasing terminal: %v", releaseErr)
	// 		}
	// 		return err
	// 	}
	// 	err = spinner.ReleaseTerminal()
	// 	if err != nil {
	// 		log.Printf("Could not release terminal: %v", err)
	// 		return err
	// 	}
	// 	fmt.Println(contents)
	//
	// 	return nil
	// }
	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
	sel := &multiSelect.Selection{
		Choices: make(map[string]bool),
	}
	optins := []string{"one", "two", "three", "four", "five"}
	ex := false
	t := tea.NewProgram(multiSelect.InitialModelMultiSelect(optins, sel, "Select", &ex))

	if _, err := t.Run(); err != nil {
		log.Fatal(err)
	}
	for i, v := range sel.Choices {
		fmt.Printf("Selected %s : is it  = %t\n", i, v)
	}
}
