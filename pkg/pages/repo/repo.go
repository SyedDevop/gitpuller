package repo

import (
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/SyedDevop/gitpuller/pkg/ui/statusbar"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type state int

// GoBackMsg is a message to go back to the previous view.
type GoBackMsg struct{}

const (
	loadingState state = iota
	readyState
)

type ReposPage struct {
	*list.Model
	statusbar *statusbar.Model
	common    common.Common
	spinner   spinner.Model

	cursor int
	state  state
}

func NewReposPage(com common.Common) *ReposPage {
	sd := statusbar.New(com)
	s := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(com.Styles.Spinner))
	repos := &ReposPage{
		statusbar: sd,
		common:    com,
		spinner:   s,
		state:     loadingState,
	}
	return repos
}

func (r *ReposPage) getMargins() (int, int) {
	hh := lipgloss.Height(r.headerView())
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		hh +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight()
	return 0, hm
}

func (r *ReposPage) SetSize(width, height int) {
	r.common.SetSize(width, height)
	_, hm := r.getMargins()
	r.statusbar.SetSize(width, height-hm)
}

// ShortHelp implements help.KeyMap.
func (r *ReposPage) ShortHelp() []key.Binding {
	k := r.KeyMap
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

func (r *ReposPage) FullHelp() [][]key.Binding {
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
	k := r.KeyMap
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

func (r *ReposPage) Title() string {
	return "Repos"
}

func (r *ReposPage) Render() string {
	return "Repo"
}
