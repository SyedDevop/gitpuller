package log

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func LogFile(path *string) *os.File {
	if path == nil || *path == "" {
		*path = "debug.log"
	}
	logF, err := tea.LogToFile(*path, "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	return logF
}
