For a CLI file downloader, you might want to structure your models and stacks to manage navigation through a directory tree and track download states. Here's a simplified example:

1. **DirectoryModel**: To represent the current state of the directory navigation.

   - Fields: `CurrentPath`, `Files`, `SelectedIndex`
   - Methods: `NavigateToSubDir()`, `NavigateToParentDir()`, `SelectFile()`

2. **DownloadModel**: To represent the state of downloads.

   - Fields: `DownloadsInProgress`, `CompletedDownloads`
   - Methods: `StartDownload()`, `CheckDownloadStatus()`

3. **NavigationStack**: To remember the path history.

   - Push the current directory to the stack before navigating to a subdirectory.
   - Pop from the stack to return to the previous directory.

4. **DownloadStack**: To track the download history or queue.
   - Push file details to the stack when a download starts.
   - Pop or query the stack to get the status of downloads.

Here's how it might work:

- The user navigates directories using `DirectoryModel`. When entering a subdirectory, the current path is pushed onto `NavigationStack`. When going back, the path is popped from `NavigationStack`.
- When the user selects a file for download, `DownloadModel` initiates the download and tracks its progress. The download details are pushed onto `DownloadStack`.

This structure allows you to separate concerns: navigation and download management, while efficiently tracking user actions and download statuses.

Certainly! Below is a simple example of how you might define a `DirectoryModel` and a `NavigationStack` for your CLI file downloader in Go:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// DirectoryModel represents the current state of the directory navigation.
type DirectoryModel struct {
	CurrentPath   string
	Files         []os.DirEntry
	SelectedIndex int
}

// NavigateToSubDir changes the current path to the selected subdirectory.
func (dm *DirectoryModel) NavigateToSubDir() error {
	if dm.SelectedIndex < 0 || dm.SelectedIndex >= len(dm.Files) {
		return fmt.Errorf("invalid selection")
	}

	selectedEntry := dm.Files[dm.SelectedIndex]
	if !selectedEntry.IsDir() {
		return fmt.Errorf("selected entry is not a directory")
	}

	name, err := selectedEntry.Name()
	if err != nil {
		return err
	}

	dm.CurrentPath = filepath.Join(dm.CurrentPath, name)
	dm.SelectedIndex = 0 // reset selection in the new directory
	return nil
}

// NavigateToParentDir changes the current path to the parent directory.
func (dm *DirectoryModel) NavigateToParentDir() {
	dm.CurrentPath = filepath.Dir(dm.CurrentPath)
	dm.SelectedIndex = 0 // reset selection in the parent directory
}

// NavigationStack is a stack of strings to keep track of directory paths.
type NavigationStack struct {
	Stack []string
}

// Push adds a new path to the stack.
func (ns *NavigationStack) Push(path string) {
	ns.Stack = append(ns.Stack, path)
}

// Pop removes and returns the top path from the stack.
func (ns *NavigationStack) Pop() (string, bool) {
	if len(ns.Stack) == 0 {
		return "", false
	}
	index := len(ns.Stack) - 1
	path := ns.Stack[index]
	ns.Stack = ns.Stack[:index]
	return path, true
}

func main() {
	// Example usage
	dirModel := DirectoryModel{
		CurrentPath:   "/path/to/start/directory",
		SelectedIndex: 0,
	}

	navStack := NavigationStack{}

	// Example navigation
	navStack.Push(dirModel.CurrentPath) // Save current path before changing
	dirModel.NavigateToSubDir()         // Navigate to a subdirectory

	// Navigate back to the previous directory
	if previousPath, ok := navStack.Pop(); ok {
		dirModel.CurrentPath = previousPath
		dirModel.SelectedIndex = 0
	}
}
```

This code provides a basic structure for navigating a file system in a CLI environment. The `DirectoryModel` keeps track of the current directory and the selected index for navigating subdirectories. The `NavigationStack` is used to remember previous directories as you navigate deeper into the file system, allowing you to backtrack correctly.
