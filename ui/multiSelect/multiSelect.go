// Package multiSelect provides functions that
// help define and draw a multi-select step
package multiSelect

import (
	"fmt"
	"strings"

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

// A Selection represents a choice made in a multiSelect step
//
//	type Selection struct {
//		Choices []types.Repo
//	}
//
//	func (s *Selection) Update(repo types.Repo) {
//		s.Choices = append(s.Choices, repo) // *(s.Choices)
//	}
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

// TODO : Check if code could be reduce.
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
			// This is sub folder root and Path.
			IsRoot, path := util.GetParentPath(m.contentTree.CurPath)
			// base file path of the hole repo
			isBasePath := m.contentTree.RootPath == m.contentTree.CurPath
			if IsRoot {
				if isBasePath {
					return m, nil
				}
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
			return m, tea.Quit
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
