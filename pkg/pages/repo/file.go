package repo

import (
	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type File struct {
	git       *gituser.Git
	items     []git.TreeElement
	common    common.Common
	TreeState StateTree
	cursor    int
}

func NewFile(com common.Common) *File {
	return &File{
		common: com,
	}
}

func (f *File) Init() tea.Cmd                           { return nil }
func (f *File) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return nil, nil }
func (f *File) View() string                            { return "From Files" }
func (f *File) SetSize(w, h int) {
	f.common.SetSize(w, h)
}

func (f *File) ShortHelp() []key.Binding {
	selectItem, open, conform := f.Keys()
	return []key.Binding{
		selectItem,
		open,
		conform,
		f.common.KeyMap.Quit,
		f.common.KeyMap.Help,
	}
}

func (f *File) FullHelp() [][]key.Binding {
	selectItem, open, conform := f.Keys()
	b := [][]key.Binding{
		{
			selectItem,
			open,
			conform,
		},
		{
			f.common.KeyMap.SelectItem,
			f.common.KeyMap.BackItem,
			f.common.KeyMap.Select,
		},
		{
			f.common.KeyMap.Home,
			f.common.KeyMap.Quit,
			f.common.KeyMap.Help,
		},
	}
	return b
}

func (f *File) Keys() (open, selectItem, conform key.Binding) {
	open = key.NewBinding(
		key.WithKeys(
			"enter",
		),
		key.WithHelp(
			"enter",
			"Open",
		),
	)

	selectItem = key.NewBinding(
		key.WithKeys(
			"space",
		),
		key.WithHelp(
			"space",
			"Select",
		),
	)

	conform = key.NewBinding(
		key.WithKeys(
			"y",
		),
		key.WithHelp(
			"y",
			"Conform the Selected Items",
		),
	)
	return
}
func (f *File) TabTitle() string { return "Files" }
