package repo

import (
	"fmt"
	"strings"

	"github.com/SyedDevop/gitpuller/pkg/assert"
	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/statusbar"
	"github.com/SyedDevop/gitpuller/pkg/ui/tabs"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type RepoFilePathMsg struct{}

type state int

const (
	loadingState state = iota
	readyState
	errorState
)

type Pane interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	SetSize(width, height int)
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
	TabTitle() string
	Reset() tea.Cmd
	StatusBarValue() string
	StatusBarInfo() string
}

type RepoPage struct {
	err          error
	git          *gituser.Git
	tabs         *tabs.Tabs
	SelectedRepo *gituser.UserRepos
	statusbar    *statusbar.Model
	repoUrl      string
	ref          string
	panes        []Pane
	common       common.Common
	spinner      spinner.Model
	activeTab    int
	state        state
}

func NewPane(com common.Common) []Pane {
	pane := []Pane{NewFile(com)}
	return pane
}

func NewRepoPage(com common.Common, gitObject *gituser.Git) *RepoPage {
	s := spinner.New(spinner.WithSpinner(spinner.Points), spinner.WithStyle(com.Styles.Spinner))
	panes := NewPane(com)
	ts := make([]string, 0)
	for _, c := range panes {
		ts = append(ts, c.TabTitle())
	}
	tb := tabs.New(com, ts)
	sb := statusbar.New(com)
	// FIX: if SelectedRepo is nil Check if repo name in context else panic.
	repos := &RepoPage{
		common:       com,
		spinner:      s,
		state:        loadingState,
		git:          gitObject,
		activeTab:    0,
		tabs:         tb,
		SelectedRepo: nil,
		panes:        panes,
		statusbar:    sb,
		ref:          "main",
	}

	return repos
}

func (r *RepoPage) getRepo(url string) tea.Cmd {
	return func() tea.Msg {
		repos, err := r.git.Repo.GetTree(url)
		if err != nil {
			return err
		}
		return repos
	}
}

func (r *RepoPage) Init() tea.Cmd {
	r.state = loadingState

	assert.Assert(r.SelectedRepo != nil, "RepoPage#GetRepo SelectedRepo cant bee nil")
	url := git.RepoUrl(r.SelectedRepo.FullName, false)
	if filePane, ok := r.panes[r.activeTab].(*File); ok {
		filePane.SetStatePath(r.SelectedRepo.Name)
	}

	if r.state == errorState {
		return func() tea.Msg {
			return r.err
		}
	}
	return tea.Batch(
		r.spinner.Tick,
		r.tabs.Init(),
		r.getRepo(url),
		r.statusbar.Init(),
	)
}

// Update implements tea.Model.
func (r *RepoPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.common.Logger.Debugf("list Msg from :%T", msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case []git.TreeElement:
		r.state = readyState
	case ReFetchRepo:
		r.state = loadingState
		cmds = append(cmds,
			r.spinner.Tick,
			r.getRepo(git.CheckDomain(string(msg))),
		)
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

	active := r.panes[r.activeTab]
	p, cmd := active.Update(msg)
	r.panes[r.activeTab] = p.(Pane)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Update the status bar on these events
	// Must come after we've updated the active tab
	switch msg.(type) {
	case ReFetchRepo, tea.KeyMsg, []git.TreeElement:
		r.setStatusBarInfo()
	}
	s, cmd := r.statusbar.Update(msg)
	r.statusbar = s.(*statusbar.Model)
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
	var statusbar string
	switch r.state {
	case loadingState:
		main = fmt.Sprintf("%s loading…", r.spinner.View())
	case readyState:
		ss := r.common.Renderer.NewStyle().
			Width(r.common.Width - wm).
			Height(r.common.Height - hm)
		main = ss.Render(r.panes[r.activeTab].View())
		statusbar = r.statusbar.View()
	}
	main = mainStyle.Render(main)
	view := lipgloss.JoinVertical(lipgloss.Top, r.headerView(), r.tabs.View(), main, statusbar)

	return s.Render(view)
}

func (r *RepoPage) Reset() tea.Cmd {
	re := r.panes[r.activeTab].Reset()
	r.activeTab = 0
	return tea.Batch(re, r.Init())
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

func (r *RepoPage) setStatusBarInfo() {
	if r.SelectedRepo == nil {
		return
	}

	active := r.panes[r.activeTab]
	key := r.SelectedRepo.Name
	value := active.StatusBarValue()
	info := active.StatusBarInfo()
	ref := "*"
	if r.ref != "" {
		ref = r.ref
	}

	r.statusbar.SetStatus(key, value, info, ref)
}

func (r *RepoPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView())
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		hh +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight() + 1
	return 0, hm
}

func (r *RepoPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	_, hm := r.getMargins()
	r.tabs.SetSize(width, height-hm)
	r.statusbar.SetSize(width, height-hm)
	for _, p := range r.panes {
		p.SetSize(width, height-hm)
	}
}

func (r *RepoPage) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	back := r.common.KeyMap.Back
	back.SetHelp("esc", "back to menu")
	tab := r.common.KeyMap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, back)
	b = append(b, tab)
	return b
}

// ShortHelp implements help.KeyMap.
func (r *RepoPage) ShortHelp() []key.Binding {
	if r.state == loadingState {
		return []key.Binding{}
	}
	b := r.commonHelp()
	b = append(b, r.panes[r.activeTab].ShortHelp()...)
	return b
}

func (r *RepoPage) FullHelp() [][]key.Binding {
	if r.state == loadingState {
		return [][]key.Binding{}
	}
	b := make([][]key.Binding, 0)
	b = append(b, r.commonHelp())
	b = append(b, r.panes[r.activeTab].FullHelp()...)
	return b
}

func (r *RepoPage) SetRepoUrl(url string) {
	r.repoUrl = url
	r.common.SetRepoUrl(url)
}

func (r *RepoPage) SetRepo(repo *gituser.UserRepos) {
	r.SelectedRepo = repo
	r.ref = repo.DefaultBranch
}
