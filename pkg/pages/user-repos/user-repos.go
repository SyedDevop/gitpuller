package userrepos

import (
	"errors"
	"fmt"

	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type state int

// GoBackMsg is a message to go back to the previous view.
type GoBackMsg struct{}

// RepoSelectedMsg is a message that is sent when a repo is selected.
type RepoSelectedMsg string

const (
	loadingState state = iota
	readyState
	errorState
)

type UserReposPage struct {
	list    list.Model
	err     error
	git     *gituser.GitUser
	common  common.Common
	spinner spinner.Model
	cursor  int
	state   state
}

func NewReposPage(com common.Common) *UserReposPage {
	s := spinner.New(spinner.WithSpinner(spinner.Points), spinner.WithStyle(com.Styles.Spinner))
	list := list.New([]list.Item{}, NewItemDelegate(&com), com.Width, com.Height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.DisableQuitKeybindings()

	per := 20
	page := 1

	g := gituser.NewGitUser()

	repos := &UserReposPage{
		common:  com,
		spinner: s,
		state:   loadingState,
		list:    list,
		git:     g,
	}
	repos.list.SetSize(com.Width, com.Height)

	// TODO: extract this to repos method.
	// this is set from --user flag
	if viper.IsSet("user") {
		name := viper.GetString("user")
		repos.git.SetUserName(name)
		repos.git.Repos.SetNextLink(git.AddPaginationParams(git.UserReposURL(name), &per, &page))
	} else if viper.IsSet("token") {
		name := viper.GetString("userName")
		repos.git.SetUserName(name)
		repos.git.Repos.SetNextLink(git.AddPaginationParams(git.AuthReposURL(), &per, &page))
	} else if viper.IsSet("userName") {
		// this is set from config file
		name := viper.GetString("userName")
		repos.git.SetUserName(name)
		repos.git.Repos.SetNextLink(git.AddPaginationParams(git.UserReposURL(name), &per, &page))
	} else {
		repos.state = errorState
		repos.err = errors.New("please provide a user name and token in the config file, or use the -u/--user flag to specify a user name for this session")
	}

	// FIX: this happens in windows only
	// repos.err = errors.New("Fix escp codes")
	return repos
}

func (r *UserReposPage) getUserRepos() tea.Cmd {
	// repos
	return func() tea.Msg {
		repos, err := r.git.Repos.Next()
		if err != nil {
			return err
		}
		return repos
	}
}

func (r *UserReposPage) Init() tea.Cmd {
	if r.state == errorState {
		return func() tea.Msg {
			return r.err
		}
	}
	return tea.Batch(
		r.spinner.Tick,
		r.getUserRepos(),
	)
}

// Update implements tea.Model.
func (r *UserReposPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.common.Logger.Debugf("list Msg from :%T", msg)
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case []gituser.UserRepos:
		r.common.Logger.Debugf("Got Msg from :%T\n and the len: %d is ", msg, len(msg))
		if len(msg) > 0 {
			r.git.SetUserUrl(msg[0].Owner.HTMLURL)
		}
		priviesItems := r.list.Items()
		newItems := make([]list.Item, len(msg))
		for i, v := range msg {
			newItems[i] = v
		}
		priviesItems = append(priviesItems, newItems...)
		r.list.SetItems(priviesItems)
		r.state = readyState
	case tea.KeyMsg:
		filterState := r.list.FilterState()
		switch {
		case key.Matches(msg, r.common.KeyMap.Select):
			if filterState != list.Filtering {
				cmds = append(cmds, r.SelectRepoCmd)
			}
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

	l, cmd := r.list.Update(msg)
	r.list = l
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	if r.state != loadingState && !r.git.Repos.ItraterDone {
		if l.Paginator.OnLastPage() {
			r.state = loadingState
			return r, r.getUserRepos()
		}
		if l.SettingFilter() {
			if len(l.VisibleItems()) == 0 {
				r.state = loadingState
				r.list.ResetFilter()
				return r, r.getUserRepos()
			}
		}
	}

	return r, tea.Batch(cmds...)
}

func (r *UserReposPage) View() string {
	wm, hm := r.getMargins()

	var view string
	switch r.state {
	case loadingState:
		view = fmt.Sprintf("%s loading…", r.spinner.View())
	case readyState:
		ss := r.common.Renderer.NewStyle().
			Width(r.common.Width - wm).
			Height(r.common.Height - hm)
		view = ss.Render(r.list.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, r.headerView(), view)
}

func (r *UserReposPage) SelectRepoCmd() tea.Msg {
	curItem := r.list.SelectedItem().(gituser.UserRepos)
	url := fmt.Sprintf("%s/main", curItem.TreesURL)
	r.common.SetRepoUrl(url)
	return RepoSelectedMsg(url)
}

func (r *UserReposPage) headerView() string {
	truncate := r.common.Renderer.NewStyle().MaxWidth(r.common.Width)
	header := r.git.Name()
	header = r.common.Styles.RepoSelector.User.HeaderName.Render(header)
	style := r.common.Styles.RepoSelector.User.Header.Copy().Width(r.common.Width)
	return style.Render(
		truncate.Render(header),
	)
}

func (r *UserReposPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView()) +
		r.common.Styles.RepoSelector.User.Header.GetVerticalFrameSize()
	return 0, hh
}

func (r *UserReposPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	wm, hm := r.getMargins()
	r.list.SetSize(width-wm, height-hm)
}

// ShortHelp implements help.KeyMap.
func (r *UserReposPage) ShortHelp() []key.Binding {
	k := r.list.KeyMap
	return []key.Binding{
		r.common.KeyMap.SelectItem,
		r.common.KeyMap.BackItem,
		r.common.KeyMap.Select,
		k.CursorUp,
		k.CursorDown,
		r.common.KeyMap.Quit,
		r.common.KeyMap.Help,
	}
}

func (r *UserReposPage) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	copyKey := r.common.KeyMap.Copy
	actionKeys := []key.Binding{
		copyKey,
	}
	copyKey.SetHelp("c", "copy name")
	k := r.list.KeyMap
	b = append(b, [][]key.Binding{
		{
			r.common.KeyMap.SelectItem,
			r.common.KeyMap.BackItem,
			r.common.KeyMap.Select,
		},
		{
			k.CursorUp,
			k.CursorDown,
			k.NextPage,
			k.PrevPage,
		},
		{
			k.GoToStart,
			k.GoToEnd,
		},
	}...)
	return append(b, actionKeys)
}
