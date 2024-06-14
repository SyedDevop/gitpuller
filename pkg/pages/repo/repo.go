package repo

import (
	"fmt"

	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	loadingState state = iota
	readyState
	errorState
)

type RepoPage struct {
	err     error
	git     *gituser.GitUser
	repoUrl string
	common  common.Common
	spinner spinner.Model
	state   state
}

func NewRepoPage(com common.Common) *RepoPage {
	s := spinner.New(spinner.WithSpinner(spinner.Points), spinner.WithStyle(com.Styles.Spinner))

	g := gituser.NewGitUser()

	g.SetUserName("Don't Know")
	repos := &RepoPage{
		common:  com,
		spinner: s,
		state:   loadingState,
		git:     g,
	}

	return repos
}

// func (r *RepoPage) getRepo() tea.Cmd {
// 	// repos
// 	return func() tea.Msg {
// 		repos, err := r.git.Repos.Next()
// 		if err != nil {
// 			return err
// 		}
// 		return repos
// 	}
// }

func (r *RepoPage) Init() tea.Cmd {
	r.state = readyState
	if r.state == errorState {
		return func() tea.Msg {
			return r.err
		}
	}
	return tea.Batch(
		r.spinner.Tick,
		// r.getRepo(),
	)
}

// Update implements tea.Model.
func (r *RepoPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.common.Logger.Debugf("list Msg from :%T", msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case []gituser.UserRepos:
		r.common.Logger.Debugf("Got Msg from :%T\n and the len: %d is ", msg, len(msg))
		r.state = readyState
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, r.common.KeyMap.Select):
		}
	case spinner.TickMsg:
		if r.state == loadingState && r.spinner.ID() == msg.ID {
			sp, cmd := r.spinner.Update(msg)
			r.spinner = sp
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return r, tea.Batch(cmds...)
}

func (r *RepoPage) View() string {
	wm, hm := r.getMargins()

	var view string
	switch r.state {
	case loadingState:
		view = fmt.Sprintf("%s loading…", r.spinner.View())
	case readyState:
		ss := r.common.Renderer.NewStyle().
			Width(r.common.Width - wm).
			Height(r.common.Height - hm)
		url := r.common.GetRepoUrl()

		view = ss.Render(r.common.Renderer.NewStyle().
			Foreground(lipgloss.Color("#5fd7ff")).Render(url),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, r.headerView(), view)
}

func (r *RepoPage) headerView() string {
	truncate := r.common.Renderer.NewStyle().MaxWidth(r.common.Width)
	header := r.git.Name()
	header = r.common.Styles.RepoSelector.User.HeaderName.Render(header)
	style := r.common.Styles.RepoSelector.User.Header.Copy().Width(r.common.Width)
	return style.Render(
		truncate.Render(header),
	)
}

func (r *RepoPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView()) +
		r.common.Styles.RepoSelector.User.Header.GetVerticalFrameSize()
	return 0, hh
}

func (r *RepoPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	// wm, hm := r.getMargins()
}

// ShortHelp implements help.KeyMap.
func (r *RepoPage) ShortHelp() []key.Binding {
	return []key.Binding{
		r.common.KeyMap.SelectItem,
		r.common.KeyMap.BackItem,
		r.common.KeyMap.Select,
		r.common.KeyMap.Home,
		r.common.KeyMap.Quit,
		r.common.KeyMap.Help,
	}
}

func (r *RepoPage) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	actionKeys := []key.Binding{}
	b = append(b, [][]key.Binding{
		{
			r.common.KeyMap.SelectItem,
			r.common.KeyMap.BackItem,
			r.common.KeyMap.Select,
		},
		{
			r.common.KeyMap.Home,
			r.common.KeyMap.Quit,
			r.common.KeyMap.Help,
		},
	}...)
	return append(b, actionKeys)
}

func (r *RepoPage) SetRepoUrl(url string) {
	r.repoUrl = url
	r.common.SetRepoUrl(url)
}
