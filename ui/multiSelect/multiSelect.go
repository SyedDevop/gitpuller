// Package multiSelect provides functions that
// help define and draw a multi-select step
package multiSelect

import (
	"fmt"

	types "github.com/SyedDevop/gitpuller/mytypes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// Change this
var (
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	titleStyle            = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	selectedItemStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	selectedItemDescStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	descriptionStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#40BDA3"))
	redText               = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

// A Selection represents a choice made in a multiSelect step
type Selection struct {
	Choices []types.Repo
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(repo types.Repo) {
	s.Choices = append(s.Choices, repo) // *(s.Choices)
}

// A multiSelect.model contains the data for the multiSelect step.
//
// It has the required methods that make it a bubbletea.Model
type model struct {
	selected map[int]struct{}
	choices  *Selection
	exit     *bool
	header   string
	options  []types.Repo
	cursor   int
}

func (m model) Init() tea.Cmd {
	return nil
}

// InitialModelMulti initializes a multiSelect step with
// the given data
func InitialModelMultiSelect(options []types.Repo, selection *Selection, header string, quit *bool) model {
	return model{
		options:  options,
		selected: make(map[int]struct{}),
		choices:  selection,
		header:   titleStyle.Render(header),
		exit:     quit,
	}
}

// Update is called when "things happen", it checks for
// important keystrokes to signal when to quit, change selection,
// and confirm the selection.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "y":
			for selectedKey := range m.selected {
				m.choices.Update(m.options[selectedKey])
				m.cursor = selectedKey
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View is called to draw the multiSelect step
func (m model) View() string {
	s := m.header + "\n\n"

	for i, option := range m.options {
		description := fmt.Sprintf("Type: %s Size: %s", option.Type, humanize.Bytes(option.Size))
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			option.Name = selectedItemStyle.Render(option.Name)
			description = selectedItemDescStyle.Render(description)
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = focusedStyle.Render("*")
		}

		title := focusedStyle.Render(option.Name)
		des := focusedStyle.Render(description)

		s += fmt.Sprintf("%s [%s] %s\n%s\n\n", cursor, checked, title, des)
	}

	s += fmt.Sprintf("Press %s to confirm choice. (%s to quit) \n", focusedStyle.Render("y"), redText.Render("q/ctrl+c"))
	return s
}
