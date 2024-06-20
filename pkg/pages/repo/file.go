package repo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/SyedDevop/gitpuller/cmd/util"
	"github.com/SyedDevop/gitpuller/pkg/git"
	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var (
	titleStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	redText           = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	fileType          = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Width(4)
	fileSize          = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(8).Align(lipgloss.Right)
	fileMode          = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	Directory         = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	file              = lipgloss.NewStyle()
	pathStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(0, 1, 0)
)

type File struct {
	git       *gituser.Git
	TreeState *StateTree
	keyMap    *FileKeyMap
	items     []git.TreeElement
	common    common.Common
	cursor    int
}
type ReFetchRepo string

func NewStateTree() *StateTree {
	return &StateTree{
		Tree:         make(map[string]*Node),
		SelectedRepo: make(map[string][]git.TreeElement),
		FolderRepo:   make([]git.TreeElement, 0),
		RootPath:     "",
		CurPath:      "",
	}
}

func NewFile(com common.Common) *File {
	return &File{
		common:    com,
		items:     make([]git.TreeElement, 0),
		TreeState: NewStateTree(),
		keyMap:    NewFileKeyMap(),
		cursor:    0,
	}
}

func (f *File) Init() tea.Cmd { return nil }
func (f *File) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			f.common.Logger.Debug(f.TreeState)
		}
		switch {
		case key.Matches(msg, f.common.KeyMap.Up):
			f.cursorUp()
		case key.Matches(msg, f.common.KeyMap.Down):
			f.cursorDown()
		case key.Matches(msg, f.keyMap.SelectItem):
			f.UpdateSelected()
		case key.Matches(msg, f.keyMap.SelectAll):
			f.SetAllSelected()
		case key.Matches(msg, f.keyMap.DSelectAll):
			f.RessetAllSelected()
		case key.Matches(msg, f.keyMap.Conform):
			f.TreeState.AppendSelected()
			f.common.Logger.Debug("Selected Items", "file", f.TreeState.FolderRepo, "Dir", f.TreeState.SelectedRepo)
		case key.Matches(msg, f.keyMap.UpDir):
			f.GoUpADir()
		case key.Matches(msg, f.keyMap.Open):
			if f.items[f.cursor].IsTree() {
				curNode := f.items[f.cursor]
				f.cursor = 0
				newPath := filepath.Join(f.TreeState.CurPath, curNode.Path)
				f.TreeState.CurPath = newPath
				if chachedNode, exists := f.TreeState.Tree[newPath]; exists {
					f.items = chachedNode.Repo
				} else {
					cmds = append(cmds, func() tea.Msg { return ReFetchRepo(*curNode.URL) })
				}
			} else {
				f.UpdateSelected()
			}
		}

	case []git.TreeElement:
		f.items = msg
		f.TreeState.Tree[f.TreeState.CurPath] = &Node{
			SelectedRepo: make(map[int]struct{}),
			Repo:         msg,
		}
	}
	return f, tea.Batch(cmds...)
}

func (f *File) View() string {
	ss := f.common.Styles.Repo.Base.Copy().
		Width(f.common.Width).
		Height(f.common.Height)
	mainStyle := f.common.Styles.Repo.Body.Copy().
		Height(f.common.Height)

	var s strings.Builder

	for i, option := range f.items {
		size := humanize.Bytes(0)
		if option.Size != nil {
			size = humanize.Bytes(uint64(*option.Size))
		}
		fsSize := fileSize.Render(size)
		fsType := fileType.Render(option.TreeType())
		fsMode := fileMode.Render(git.ToOSFileMode(option.Mode).String())
		description := fmt.Sprintf("%s %s %s", fsMode, fsType, fsSize)

		cursor := " "
		if f.cursor == i {
			cursor = focusedStyle.Render(">")
			option.Path = focusedStyle.Render(option.Path)
			description = focusedStyle.Render(description)
		}

		checked := " "
		if node, ok := f.TreeState.Tree[f.TreeState.CurPath]; ok {
			if _, ok := node.SelectedRepo[i]; ok {
				checked = selectedItemStyle.Render("*")
				option.Path = selectedItemStyle.Render(option.Path)
				description = selectedItemStyle.Render(description)
			}
		}

		option.Path = file.Render(option.Path)
		if option.IsTree() {
			option.Path = Directory.Render(option.Path)
		}

		s.WriteString(fmt.Sprintf("%s %s %s %s\n", cursor, checked, description, option.Path))
	}
	return ss.Render(mainStyle.Render(s.String()))
}

func (f *File) Reset() tea.Cmd {
	f.TreeState = NewStateTree()
	f.cursor = 0
	return nil
}

func (f *File) GoUpADir() {
	isBasePath := f.TreeState.RootPath == f.TreeState.CurPath
	if !isBasePath {
		_, parentPath := util.GetParentPath(f.TreeState.CurPath)
		if cachedNode, exists := f.TreeState.Tree[parentPath]; exists {
			f.items = cachedNode.Repo
			f.TreeState.CurPath = parentPath
		}
	}
}

func (f *File) UpdateSelected() {
	f.TreeState.UpdateTreesSelected(f.cursor)
}

func (f *File) SetAllSelected() {
	f.TreeState.SelectAllCurTreeRepo()
}

func (f *File) RessetAllSelected() {
	f.TreeState.RemoveAllCurTreeRepo()
}

func (f *File) SetStatePath(path string) {
	f.TreeState.CurPath = path
	f.TreeState.RootPath = path
}

func (f *File) SetSize(w, h int) {
	f.common.SetSize(w, h)
}

func (f *File) ShortHelp() []key.Binding {
	return []key.Binding{
		f.keyMap.SelectItem,
		f.keyMap.Open,
		f.keyMap.Conform,
		f.keyMap.SelectAll,
		f.keyMap.DSelectAll,
		f.common.KeyMap.Quit,
		f.common.KeyMap.Help,
	}
}

func (f *File) FullHelp() [][]key.Binding {
	b := [][]key.Binding{
		{
			f.keyMap.SelectItem,
			f.keyMap.Open,
			f.keyMap.Conform,
		},
		{
			f.keyMap.SelectAll,
			f.keyMap.DSelectAll,
			f.common.KeyMap.UpDown,
		},
		{
			f.common.KeyMap.Home,
			f.common.KeyMap.Quit,
			f.common.KeyMap.Help,
		},
	}
	return b
}

func (f *File) cursorUp() {
	if f.cursor == 0 {
		return
	}
	f.cursor--
}

func (f *File) cursorDown() {
	if f.cursor >= len(f.items)-1 {
		return
	}
	f.cursor++
}

func (f *File) TabTitle() string { return "Files" }
