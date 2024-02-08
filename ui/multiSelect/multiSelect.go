// Package multiSelect provides functions that
// help define and draw a multi-select step
package multiSelect

import (
	"fmt"

	"github.com/SyedDevop/gitpuller/api"
	types "github.com/SyedDevop/gitpuller/mytypes"
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
)

type (
	multiSelectMsg string
	errMess        struct{ error }
)

func (e errMess) Error() string { return e.error.Error() }

// A Selection represents a choice made in a multiSelect step
type Selection struct {
	Choices []types.Repo
}

type TreeData struct {
	SelectedRepo map[int]struct{}
	Repo         []types.Repo
}

type ContentTree struct {
	Tree    map[string]TreeData
	CurPath string
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(repo types.Repo) {
	s.Choices = append(s.Choices, repo) // *(s.Choices)
}

type Fetch struct {
	Err       error
	Clint     *api.Clint
	FetchMess string
	PathRoute string
	Repo      []types.Repo
	FethDone  bool
}

// A multiSelect.Model contains the data for the multiSelect step.
//
// It has the required methods that make it a bubbletea.Model
type Model struct {
	selected    map[int]struct{}
	choices     *Selection
	exit        *bool
	fetch       *Fetch
	contentTree *ContentTree
	header      string
	options     []types.Repo
	spinner     spinner.Model
	cursor      int
}

// InitialModelMulti initializes a multiSelect step with
// the given data
func InitialModelMultiSelect(clintFetch *Fetch, selection *Selection, conTree *ContentTree, header string, quit *bool) Model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return Model{
		options:     make([]types.Repo, 0),
		selected:    make(map[int]struct{}),
		choices:     selection,
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
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "backspace", "b":
			// if len(m.contentTree.PathRoute) <= 1 {
			// 	return m, tea.Println("Last")
			// }
			// m.contentTree.PathRoute = m.contentTree.PathRoute[:len(m.contentTree.PathRoute)-1]
			// preRepo := m.contentTree.PathRoute[len(m.contentTree.PathRoute)-1]
			// data, ok := m.contentTree.Tree[preRepo]
			// if ok {
			// 	m.selected = data.SelectedRepo
			// 	m.options = data.Repo
			// }
			return m, tea.Batch(tea.Println(m.contentTree.CurPath))

		case "enter":
			if m.options[m.cursor].Type != "dir" {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			} else {
				curDir := m.options[m.cursor]
				data, ok := m.contentTree.Tree[curDir.Path]
				if ok {
					m.selected = data.SelectedRepo
					m.options = data.Repo
					return m, tea.Batch(tea.Println("Cached"))
				}
				m.fetch.FethDone = false
				m.fetch.Clint.GitRepoUrl = curDir.URL
				m.fetch.PathRoute = curDir.Path
				m.cursor = 0
				return m, tea.Batch(m.fetch.fetchContent)
			}

		case "a", "A":
			if len(m.options) > len(m.selected) {
				for i := 0; i < len(m.options); i++ {
					m.selected[i] = struct{}{}
				}
			}
		case "d", "D":
			for i := 0; i < len(m.options); i++ {
				delete(m.selected, i)
			}
		case "y":
			for selectedKey := range m.selected {
				m.choices.Update(m.options[selectedKey])
				m.cursor = selectedKey
			}
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case multiSelectMsg:
		m.fetch.FethDone = true
		m.options = m.fetch.Repo

		for k := range m.selected {
			delete(m.selected, k)
		}

		m.contentTree.CurPath = m.fetch.PathRoute
		// m.contentTree.PathRoute = append(m.contentTree.PathRoute, m.fetch.PathRoute)
		m.contentTree.Tree[m.fetch.PathRoute] = TreeData{
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

// View is called to draw the multiSelect step
func (m Model) View() string {
	s := m.header + "\n\n"
	if !m.fetch.FethDone {
		s += fmt.Sprintf("%s %s... Press 'q' to quit", m.spinner.View(), m.fetch.FetchMess)
		return s
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
		if _, ok := m.selected[i]; ok {
			checked = selectedItemStyle.Render("*")
			option.Name = selectedItemStyle.Render(option.Name)
			description = selectedItemStyle.Render(description)
		}

		option.Name = File.Render(option.Name)
		if option.Type != "file" {
			option.Name = Directory.Render(option.Name)
		}

		// title := focusedStyle.Render(option.Name)

		s += fmt.Sprintf("%s %s %s %s\n", cursor, checked, description, option.Name)
	}

	s += fmt.Sprintf("\nPress %s to confirm choice. (%s to quit) \n", selectedItemStyle.Render("y"), redText.Render("q/ctrl+c"))
	return s
}
