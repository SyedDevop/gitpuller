// Package multiSelect provides functions that
// help define and draw a multi-select step
package multiSelect

import (
	"fmt"
	"strings"
	"sync"

	types "github.com/SyedDevop/gitpuller/mytypes"
	"github.com/SyedDevop/gitpuller/util"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var (
	titleStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	redText           = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	fileType          = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Width(4)
	fileSize          = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(8).Align(lipgloss.Right)
	Directory         = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	File              = lipgloss.NewStyle()
	pathStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(0, 1, 0)
)

type (
	multiSelectMsg string
	errMess        struct{ error }
)

func (e errMess) Error() string { return e.error.Error() }

type Model struct {
	exit        *bool
	fetch       *Fetch
	contentTree *ContentTree
	header      string
	options     []types.Repo
	spinner     spinner.Model
	cursor      int
}

func InitialModelMultiSelect(clintFetch *Fetch, conTree *ContentTree, header string, quit *bool) Model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return Model{
		options:     make([]types.Repo, 0),
		header:      titleStyle.Render(header),
		exit:        quit,
		spinner:     s,
		fetch:       clintFetch,
		contentTree: conTree,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetch.fetchContent, m.spinner.Tick)
}

// TODO: Check if code could be reduce.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case " ":
			m.contentTree.UpdateTreesSelected(m.cursor)

		case "backspace", "b":

			// base file path of the whole repo
			isBasePath := m.contentTree.RootPath == m.contentTree.CurPath
			if isBasePath {
				return m, nil
			}

			// This is sub folder of root and Path.
			IsRoot, path := util.GetParentPath(m.contentTree.CurPath)
			if IsRoot {
				path = m.contentTree.RootPath
			}
			chachedNode, ok := m.contentTree.Tree[path]
			if ok {
				m.options = chachedNode.Repo
				m.contentTree.CurPath = path
			}
			return m, nil

		case "enter":
			if m.options[m.cursor].Type != "dir" {
				m.contentTree.UpdateTreesSelected(m.cursor)
			} else {
				curDir := m.options[m.cursor]
				chachedNode, ok := m.contentTree.Tree[curDir.Path]
				m.cursor = 0
				m.contentTree.CurPath = curDir.Path

				if ok {
					m.options = chachedNode.Repo
					return m, nil
				}
				m.fetch.FethDone = false
				m.fetch.Clint.GitRepoUrl = curDir.URL
				return m, m.fetch.fetchContent
			}

		case "a", "A":
			m.contentTree.SelectAllCurTreeRepo()
		case "d", "D":
			m.contentTree.RemoveAllCurTreeRepo()
		case "y":
			m.contentTree.AppendSelected()
			m.fetch.FethDone = false
			m.fetch.FetchMess = "Fetching Repo Files..."
			return m, FetchAllFolders(&m)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case multiSelectMsg:
		m.fetch.FethDone = true
		m.options = m.fetch.Repo

		m.contentTree.Tree[m.contentTree.CurPath] = &Node{
			SelectedRepo: make(map[int]struct{}),
			Repo:         m.options,
		}
		return m, nil
	case errMess:
		m.fetch.Err = msg
		return m, tea.Quit
	}

	return m, nil
}

func FetchAllFolders(model *Model) tea.Cmd {
	return func() tea.Msg {
		wg := sync.WaitGroup{}
		list := model.contentTree.FolderRepo
		wg.Add(len(list))
		errChan := make(chan error)
		for _, repo := range list {
			go func(repo types.Repo) {
				defer wg.Done()
				allRepos, err := FetchRepoFiles(repo.URL, model.fetch)
				if err != nil {
					errChan <- err
				}
				model.contentTree.Mu.Lock()
				model.contentTree.SelectedRepo = append(model.contentTree.SelectedRepo, allRepos...)
				model.contentTree.Mu.Unlock()
			}(repo)
		}

		// TODO: check if this can be done in a better way
		wg.Wait()
		close(errChan)

		// TODO: Try to return error as list of errors
		for err := range errChan {
			if err != nil {
				return errMess{err}
			}
		}
		return tea.QuitMsg{}
	}
}

func FetchRepoFiles(url string, fetch *Fetch) ([]types.Repo, error) {
	var repos []types.Repo
	fetch.Clint.GitRepoUrl = url
	data, err := fetch.Clint.GetCountents()
	if err != nil {
		return nil, err
	}
	rawData := util.GetRepoFromContent(*data)

	for _, item := range rawData {
		if item.Type == "dir" {
			// Recursively fetch contents from the directory, excluding the directory itself
			newData, err := FetchRepoFiles(item.URL, fetch)
			if err != nil {
				return nil, err
			}
			repos = append(repos, newData...)
		} else {
			// Append non-directory items to the list
			repos = append(repos, item)
		}
	}
	return repos, nil
}

// View is called to draw the multiSelect step
func (m Model) View() string {
	var s strings.Builder
	currebtPath := pathStyle.Render("Current Path: (" + m.contentTree.CurPath + ")")
	s.WriteString(m.header + "\n" + currebtPath + "\n\n")
	if !m.fetch.FethDone {
		s.WriteString(fmt.Sprintf("%s %s... Press 'q' to quit", m.spinner.View(), m.fetch.FetchMess))
		return s.String()
	}

	for i, option := range m.options {
		fsSize := fileSize.Render(humanize.Bytes(option.Size))
		fsType := fileType.Render(option.Type)
		description := fmt.Sprintf("%s %s", fsType, fsSize)

		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			option.Name = focusedStyle.Render(option.Name)
			description = focusedStyle.Render(description)
		}

		checked := " "
		if _, ok := m.contentTree.Tree[m.contentTree.CurPath].SelectedRepo[i]; ok {
			checked = selectedItemStyle.Render("*")
			option.Name = selectedItemStyle.Render(option.Name)
			description = selectedItemStyle.Render(description)
		}

		option.Name = File.Render(option.Name)
		if option.Type != "file" {
			option.Name = Directory.Render(option.Name)
		}

		// title := focusedStyle.Render(option.Name)

		s.WriteString(fmt.Sprintf("%s %s %s %s\n", cursor, checked, description, option.Name))
	}

	s.WriteString(fmt.Sprintf("\nPress %s to confirm choice. (%s to quit) \n", selectedItemStyle.Render("y"), redText.Render("q/ctrl+c")))
	return s.String()
}
