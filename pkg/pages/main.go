package pages

import (
	"io"

	userrepos "github.com/SyedDevop/gitpuller/pkg/pages/user-repos"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/footer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type Page interface {
	Title() string
	Render() string
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
	Init() tea.Cmd
	View() string
	SetSize(width, height int)
	Update(msg tea.Msg) (tea.Model, tea.Cmd)

	// ProjectName() string
	// Name() string
	// Description() string
}

type Model struct {
	error       error
	footer      *footer.Footer
	pages       []Page
	common      common.Common
	currentPage int
	state       state
	showFooter  bool
}

type (
	state int
)

const (
	startState state = iota
	errorState
)

func NewPageModel(cmd *cobra.Command, fileLogger io.Writer) *Model {
	output := lipgloss.DefaultRenderer()
	ctx := cmd.Context()
	c := common.NewCommon(ctx, fileLogger, output, 0, 0)

	r := userrepos.NewReposPage(c)
	m := &Model{
		common:      c,
		currentPage: 0,
		pages: []Page{
			r,
		},
	}

	m.footer = footer.New(c, m)
	return m
}

func (m *Model) SetSize(w, h int) {
	m.common.SetSize(w, h)
	style := m.common.Styles.App.Copy()
	wm := style.GetHorizontalFrameSize()
	hm := style.GetVerticalFrameSize()
	if m.showFooter {
		hm += m.footer.Height()
	}

	m.footer.SetSize(w-wm, h-hm)
	m.pages[0].SetSize(w-wm, h-hm)
}

// ShortHelp implements help.KeyMap.
func (m Model) ShortHelp() []key.Binding {
	switch m.state {
	case errorState:
		return []key.Binding{
			m.common.KeyMap.Back,
			m.common.KeyMap.Quit,
			m.common.KeyMap.Help,
		}
	default:

		// FIX : Change to use current Page/panes help
		// return []key.Binding{
		// 	m.common.KeyMap.Back,
		// 	m.common.KeyMap.Quit,
		// 	m.common.KeyMap.Help,
		// }
		return m.pages[0].ShortHelp()
	}
}

// FullHelp implements help.KeyMap.
func (m Model) FullHelp() [][]key.Binding {
	switch m.state {
	case errorState:
		return [][]key.Binding{
			{
				m.common.KeyMap.Back,
			},
			{
				m.common.KeyMap.Quit,
				m.common.KeyMap.Help,
			},
		}
	default:
		// FIX : Change to use current Page/panes help
		// return [][]key.Binding{
		// 	{
		// 		m.common.KeyMap.Back,
		// 	},
		// 	{
		// 		m.common.KeyMap.Quit,
		// 		m.common.KeyMap.Help,
		// 	},
		// }
		return m.pages[0].FullHelp()
	}
}

func (m *Model) Init() tea.Cmd {
	// FIX : Change to use current Page/panes init
	rePage := m.pages[0]
	m.state = startState
	return tea.Batch(m.footer.Init(), rePage.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// m.common.Logger.Debugf("mgs received: %T", msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case footer.ToggleFooterMsg:
		m.footer.SetShowAll(!m.footer.ShowAll())
		m.showFooter = !m.showFooter
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.common.KeyMap.Back) && m.error != nil:
			m.error = nil
			m.state = startState
			// Always show the footer on error.
			m.showFooter = m.footer.ShowAll()
		case key.Matches(msg, m.common.KeyMap.Help):
			cmds = append(cmds, footer.ToggleFooterCmd)
		case key.Matches(msg, m.common.KeyMap.Quit):
			return m, tea.Quit
		}

	case common.ErrorMsg:
		m.error = msg
		m.state = errorState
		m.showFooter = true
	}

	// NOTE: This is how you can handle different modal
	f, cmd := m.footer.Update(msg)
	m.footer = f.(*footer.Footer)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// FIX : Change to use current Page/panes
	rePage, cmd := m.pages[0].Update(msg)
	m.pages[0] = rePage.(*userrepos.UserReposPage)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// This fixes determining the height margin of the footer.
	m.SetSize(m.common.Width, m.common.Height)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	myStyle := m.common.Styles.App.Copy()

	wm, hm := myStyle.GetHorizontalFrameSize(), myStyle.GetVerticalFrameSize()
	if m.showFooter {
		hm += m.footer.Height()
	}

	var view string
	switch m.state {
	case startState:
		// FIX : Change to use current Page/panes
		view = m.pages[0].View()
	case errorState:
		err := m.common.Styles.ErrorTitle.Render("Bummer")
		err += m.common.Styles.ErrorBody.Render(m.error.Error())
		view = m.common.Styles.Error.Copy().
			Width(m.common.Width -
				wm -
				m.common.Styles.ErrorBody.GetHorizontalFrameSize()).
			Height(m.common.Height -
				hm -
				m.common.Styles.Error.GetVerticalFrameSize()).
			Render(err)
	}
	if m.showFooter {
		view = lipgloss.JoinVertical(lipgloss.Top, view, m.footer.View())
	}

	view = lipgloss.JoinVertical(lipgloss.Top, view)
	return myStyle.Render(view)
}
