package repo

import (
	"fmt"

	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// var (
// 	titleStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
// 	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
// 	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#73DACA")).Bold(true)
// 	redText           = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
// 	fileType          = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Width(4)
// 	fileSize          = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(8).Align(lipgloss.Right)
// 	fileMode          = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	Directory         = lipgloss.NewStyle().Foreground(lipgloss.Color("#2AC3DE"))
// 	file              = lipgloss.NewStyle()
// 	pathStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true).Padding(0, 1, 0)
// )

type Branch struct {
	list   list.Model
	common common.Common
}

func NewBranch(com common.Common) *Branch {
	list := list.New([]list.Item{}, list.DefaultDelegate{}, com.Width, com.Height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.DisableQuitKeybindings()
	return &Branch{
		common: com,
		list:   list,
	}
}

func (b *Branch) Init() tea.Cmd { return nil }
func (b *Branch) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	l, cmd := b.list.Update(msg)
	b.list = l
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return b, tea.Batch(cmds...)
}

func (b *Branch) View() string {
	return b.list.View()
}

func (b *Branch) Reset() tea.Cmd {
	return nil
}

// StatusBarValue returns the status bar value.
func (b *Branch) StatusBarValue() string {
	p := b.list.SelectedItem().FilterValue()
	return p
}

// StatusBarInfo returns the status bar info.
func (b *Branch) StatusBarInfo() string {
	return fmt.Sprintf("# %d/%d", b.list.Paginator.Page, b.list.Paginator.TotalPages)
}

func (b *Branch) SetSize(w, h int) {
	b.common.SetSize(w, h)
}

func (b *Branch) ShortHelp() []key.Binding {
	return []key.Binding{
		b.common.KeyMap.Quit,
		b.common.KeyMap.Help,
	}
}

func (b *Branch) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{
		{
			b.common.KeyMap.UpDown,
		},
		{
			b.common.KeyMap.Home,
			b.common.KeyMap.Quit,
			b.common.KeyMap.Help,
		},
	}
	return kb
}

func (b *Branch) TabTitle() string { return "Branches" }
