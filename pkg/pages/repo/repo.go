package repo

import (
	"fmt"
	"strings"

	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/tabs"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type state int

const (
	loadingState state = iota
	readyState
	errorState
)

type RepoPage struct {
	err          error
	git          *gituser.Git
	tabs         *tabs.Tabs
	SelectedRepo *gituser.UserRepos
	repoUrl      string
	common       common.Common
	spinner      spinner.Model
	activeTab    int
	state        state
}

func NewRepoPage(com common.Common, gitObject *gituser.Git) *RepoPage {
	s := spinner.New(spinner.WithSpinner(spinner.Points), spinner.WithStyle(com.Styles.Spinner))

	tb := tabs.New(com, []string{"file", "branch", "tag"})

	// FIX: if SelectedRepo is nil Check if repo name in context else panic.
	repos := &RepoPage{
		common:       com,
		spinner:      s,
		state:        loadingState,
		git:          gitObject,
		activeTab:    0,
		tabs:         tb,
		SelectedRepo: nil,
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
		r.tabs.Init(),
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
	case tabs.SelectTabMsg:
		r.activeTab = int(msg)
		t, cmd := r.tabs.Update(msg)
		r.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		r.activeTab = int(msg)
	}

	t, cmd := r.tabs.Update(msg)
	r.tabs = t.(*tabs.Tabs)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return r, tea.Batch(cmds...)
}

func (r *RepoPage) View() string {
	wm, hm := r.getMargins()
	hm += r.common.Styles.Tabs.GetHeight() +
		r.common.Styles.Tabs.GetVerticalFrameSize()
	s := r.common.Styles.Repo.Base.Copy().
		Width(r.common.Width - wm).
		Height(r.common.Height - hm)
	mainStyle := r.common.Styles.Repo.Body.Copy().
		Height(r.common.Height - hm)

	var main string
	switch r.state {
	case loadingState:
		main = fmt.Sprintf("%s loading…", r.spinner.View())
	case readyState:
		ss := r.common.Renderer.NewStyle().
			Width(r.common.Width - wm).
			Height(r.common.Height - hm)
		url := fmt.Sprintf("Current Tab %s (Tab Index %d)\n%s", r.tabs.GetActiveTabName(), r.activeTab, r.common.GetRepoUrl())

		main = ss.Render(r.common.Renderer.NewStyle().
			Foreground(lipgloss.Color("#5fd7ff")).Render(url),
		)
	}
	main = mainStyle.Render(main)
	view := lipgloss.JoinVertical(lipgloss.Top, r.headerView(), r.tabs.View(), main)

	return s.Render(view)
}

func TruncateString(s string, max int) string {
	if max < 0 {
		max = 0
	}
	return truncate.StringWithTail(s, uint(max), "…")
}

func (r *RepoPage) headerView() string {
	if r.SelectedRepo == nil {
		return ""
	}
	truncate := r.common.Renderer.NewStyle().MaxWidth(r.common.Width)
	header := r.SelectedRepo.Name

	header = r.common.Styles.Repo.HeaderName.Render(header)
	desc := strings.TrimSpace(r.SelectedRepo.Description())
	header = lipgloss.JoinVertical(lipgloss.Top,
		header,
		r.common.Styles.Repo.HeaderDesc.Render(desc),
	)
	urlStyle := r.common.Styles.URLStyle.Copy().
		Width(r.common.Width - lipgloss.Width(desc) - 1).
		Align(lipgloss.Right).MarginTop(1)
	url := r.SelectedRepo.Command()
	url = TruncateString(url, r.common.Width-lipgloss.Width(desc)-1)
	url = urlStyle.Render(url)

	header = lipgloss.JoinHorizontal(lipgloss.Left, header, url)
	style := r.common.Styles.Repo.Header.Copy().Width(r.common.Width)
	return style.Render(
		truncate.Render(header),
	)
}

func (r *RepoPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView())
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		hh +
		r.common.Styles.Repo.Header.GetVerticalFrameSize()
	return 0, hm
}

func (r *RepoPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	_, hm := r.getMargins()
	r.tabs.SetSize(width, height-hm)
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

func (r *RepoPage) SetRepo(repo *gituser.UserRepos) {
	r.SelectedRepo = repo
}
