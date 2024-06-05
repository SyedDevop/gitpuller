package userrepos

import (
	"fmt"
	"strings"

	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/statusbar"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

// GoBackMsg is a message to go back to the previous view.
type GoBackMsg struct{}

const (
	loadingState state = iota
	readyState
)

type UserReposPage struct {
	list      list.Model
	statusbar *statusbar.Model
	git       *gituser.GitUser
	common    common.Common
	spinner   spinner.Model
	cursor    int
	state     state
}

func NewReposPage(com common.Common) *UserReposPage {
	sd := statusbar.New(com)
	s := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(com.Styles.Spinner))
	list := list.New([]list.Item{}, NewItemDelegate(&com), com.Width, com.Height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	g := gituser.NewGitUser("SyedDevop")
	per := 20
	page := 1

	// FIX: change to be authenticated of user based ReposLink
	// - [] if there is a user Name
	// - [] else if there is a authenticated user
	// - [] else if there is a default user
	// - [] else Panic and say any one of the option needs to be provided,
	g.Repos.SetNextLink(git.AddPaginationParams(git.AuthReposURL(), &per, &page))

	repos := &UserReposPage{
		statusbar: sd,
		common:    com,
		spinner:   s,
		state:     loadingState,
		list:      list,
		git:       g,
	}
	repos.list.SetSize(com.Width, com.Height)
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
	return tea.Batch(
		r.statusbar.Init(),
		r.spinner.Tick,
		r.getUserRepos(),
	)
}

// Update implements tea.Model.
func (r *UserReposPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case spinner.TickMsg:
		sp, cmd := r.spinner.Update(msg)
		r.spinner = sp
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	// Update the status bar on these events
	// Must come after we've updated the active tab

	l, cmd := r.list.Update(msg)
	r.list = l
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	if r.state != loadingState && !r.git.Repos.ItraterDone {
		if l.Paginator.OnLastPage() {
			r.state = loadingState
			r.page = l.Paginator.Page
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

	s, cmd := r.statusbar.Update(msg)
	r.statusbar = s.(*statusbar.Model)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return r, tea.Batch(cmds...)
}

func (r *UserReposPage) View() string {
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
		main = r.list.View()
		statusbar = r.statusbar.View()
	}

	// wordSty := r.common.Renderer.NewStyle().MaxWidth(r.common.Width)
	// word := ""
	// if r.list.FilterState() == list.Filtering {
	// 	word = fmt.Sprintf("Currently showing %d results", len(r.list.VisibleItems()))
	// }

	view := lipgloss.JoinVertical(lipgloss.Top,
		r.headerView(),
		// word,
		mainStyle.Render(main),
		statusbar,
	)
	return s.Render(view)
}

func (r *UserReposPage) headerView() string {
	truncate := r.common.Renderer.NewStyle().MaxWidth(r.common.Width)
	header := r.git.ProjectName()
	if header == "" {
		header = r.git.Name()
	}
	header = r.common.Styles.Repo.HeaderName.Render(header)
	desc := strings.TrimSpace(r.git.Description())
	if desc != "" {
		header = lipgloss.JoinVertical(lipgloss.Top,
			header,
			r.common.Styles.Repo.HeaderDesc.Render(desc),
		)
	}
	// urlStyle := r.common.Styles.URLStyle.Copy().
	// 	Width(r.common.Width - lipgloss.Width(desc) - 1).
	// 	Align(lipgloss.Right)
	// var url string
	// if cfg := r.common.Config(); cfg != nil {
	// 	url = r.common.CloneCmd(cfg.SSH.PublicURL, r.repos.Name())
	// }
	// url = common.TruncateString(url, r.common.Width-lipgloss.Width(desc)-1)
	// url = r.common.Zone.Mark(
	// 	fmt.Sprintf("%s-url", r.repos.Name()),
	// 	urlStyle.Render(url),
	// )

	header = lipgloss.JoinHorizontal(lipgloss.Left, header)

	style := r.common.Styles.Repo.Header.Copy().Width(r.common.Width)
	return style.Render(
		truncate.Render(header),
	)
}

func (r *UserReposPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView())
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		hh +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight()
	return 0, hm
}

func (r *UserReposPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	_, hm := r.getMargins()
	r.list.SetSize(width, height-hm)
	r.statusbar.SetSize(width, height-hm)
}

// ShortHelp implements help.KeyMap.
func (r *UserReposPage) ShortHelp() []key.Binding {
	k := r.list.KeyMap
	return []key.Binding{
		r.common.KeyMap.SelectItem,
		r.common.KeyMap.BackItem,
		k.CursorUp,
		k.CursorDown,
	}
	// case filesViewContent:
	// 	b := []key.Binding{
	// 		f.common.KeyMap.UpDown,
	// 		f.common.KeyMap.BackItem,
	// 	}
	// 	return b
	// default:
	// 	return []key.Binding{}
	// }
}

func (r *UserReposPage) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	copyKey := r.common.KeyMap.Copy
	actionKeys := []key.Binding{
		copyKey,
	}
	// if !f.code.UseGlamour {
	// 	actionKeys = append(actionKeys, lineNo)
	// }
	//  TODO: implement this for readMe
	// if common.IsFileMarkdown(f.currentContent.content, f.currentContent.ext) &&
	// 	!f.blameView {
	// 	actionKeys = append(actionKeys, preview)
	// }
	// switch f.activeView {
	// case filesViewFiles:
	copyKey.SetHelp("c", "copy name")
	k := r.list.KeyMap
	b = append(b, [][]key.Binding{
		{
			r.common.KeyMap.SelectItem,
			r.common.KeyMap.BackItem,
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
	// case filesViewContent:
	// 	copyKey.SetHelp("c", "copy content")
	// 	k := f.code.KeyMap
	// 	b = append(b, []key.Binding{
	// 		f.common.KeyMap.BackItem,
	// 	})
	// 	b = append(b, [][]key.Binding{
	// 		{
	// 			k.PageDown,
	// 			k.PageUp,
	// 			k.HalfPageDown,
	// 			k.HalfPageUp,
	// 		},
	// 		{
	// 			k.Down,
	// 			k.Up,
	// 			f.common.KeyMap.GotoTop,
	// 			f.common.KeyMap.GotoBottom,
	// 		},
	// 	}...)
	// }
	return append(b, actionKeys)
}

func (r *UserReposPage) Title() string {
	return "Repos"
}

func (r *UserReposPage) Render() string {
	return "Repo"
}
