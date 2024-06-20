package pages

import (
	"io"

	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/pages/repo"
	userrepos "github.com/SyedDevop/gitpuller/pkg/pages/user-repos"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/footer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type Page interface {
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

type Pane int

const (
	selectionPage Pane = iota
	repoPage
)

type Model struct {
	error       error
	footer      *footer.Footer
	pages       []Page
	common      common.Common
	currentPage Pane
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

	// TODO: Move the gituser to git module
	git := gituser.NewGitUser()
	userRposPage := userrepos.NewReposPage(c, git)
	RepoPage := repo.NewRepoPage(c, git)
	m := &Model{
		common:      c,
		currentPage: selectionPage,
		pages: []Page{
			userRposPage,
			RepoPage,
		},
	}

	m.footer = footer.New(c, m)
	m.showFooter = true
	return m
}

func (m *Model) getMargins() (wm, hm int) {
	style := m.common.Styles.App.Copy()
	// switch m.activePage {
	// case selectionPage:
	// 	hm += m.common.Styles.ServerName.GetHeight() +
	// 		m.common.Styles.ServerName.GetVerticalFrameSize()
	// case repoPage:
	// }
	wm += style.GetHorizontalFrameSize()
	hm += style.GetVerticalFrameSize()
	if m.showFooter {
		// NOTE: we don't use the footer's style to determine the margins
		// because footer.Height() is the height of the footer after applying
		// the styles.
		hm += m.footer.Height()
	}
	return
}

func (m *Model) SetSize(w, h int) {
	m.common.SetSize(w, h)
	wm, hm := m.getMargins()
	// wm := style.GetHorizontalFrameSize()
	// hm := style.GetVerticalFrameSize()
	// if m.footer.ShowAll() {
	// 	hm += m.footer.Height()
	// }

	m.footer.SetSize(w-wm, h-hm)
	m.pages[m.currentPage].SetSize(w-wm, h-hm)
}

// ShortHelp implements help.KeyMap.
func (m Model) ShortHelp() []key.Binding {
	switch m.state {
	case errorState:
		return []key.Binding{
			m.common.KeyMap.Refresh,
			m.common.KeyMap.Quit,
			m.common.KeyMap.Help,
		}
	default:
		return m.pages[m.currentPage].ShortHelp()
	}
}

// FullHelp implements help.KeyMap.
func (m Model) FullHelp() [][]key.Binding {
	switch m.state {
	case errorState:
		return [][]key.Binding{
			{
				m.common.KeyMap.Refresh,
			},
			{
				m.common.KeyMap.Quit,
				m.common.KeyMap.Help,
			},
		}
	default:
		return m.pages[m.currentPage].FullHelp()
	}
}

func (m *Model) Init() tea.Cmd {
	rePage := m.pages[m.currentPage]
	m.state = startState
	return tea.Batch(m.footer.Init(), rePage.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.common.Logger.Debugf("list Msg from :%T", msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case userrepos.RepoSelectedMsg:
		m.currentPage = repoPage
		curPage := m.pages[m.currentPage]
		if repoPage, ok := curPage.(*repo.RepoPage); ok {
			repoPage.SetRepo(&msg)
		}
		cmd := curPage.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.showFooter = false

	case footer.ToggleFooterMsg:
		m.footer.SetShowAll(!m.footer.ShowAll())
		m.showFooter = !m.showFooter
	case tea.KeyMsg:
		switch {
		// Request to go back
		case key.Matches(msg, m.common.KeyMap.Refresh) && m.error != nil:
			m.error = nil
			m.state = startState
			// Always show the footer on error.
			m.showFooter = m.footer.ShowAll()
			cmds = append(cmds, m.pages[m.currentPage].Init())
		case key.Matches(msg, m.common.KeyMap.Home):
			if m.currentPage != selectionPage {
				m.currentPage = selectionPage
				curPage := m.pages[m.currentPage]
				cmd := curPage.Init()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
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
	if !m.showFooter && m.currentPage == selectionPage {
		m.common.Logger.Debug("Showing footer", "isFooterShowing", m.showFooter)
		m.showFooter = true
	}
	f, cmd := m.footer.Update(msg)
	m.footer = f.(*footer.Footer)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// FIX : Change to use current Page/panes
	rePage, cmd := m.pages[m.currentPage].Update(msg)
	if m.currentPage == selectionPage {
		m.pages[m.currentPage] = rePage.(*userrepos.UserReposPage)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// This fixes determining the height margin of the footer.
	m.SetSize(m.common.Width, m.common.Height)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	wm, hm := m.getMargins()
	var view string
	switch m.state {
	case startState:
		// FIX : Change to use current Page/panes
		view = m.pages[m.currentPage].View()
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

	return m.common.Styles.App.Render(view)
}
