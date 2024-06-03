package userrepos

import (
	"fmt"
	"io"
	"strings"

	gituser "github.com/SyedDevop/gitpuller/pkg/git/git-user"
	"github.com/SyedDevop/gitpuller/pkg/ui/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/muesli/reflow/truncate"
)

type ItemDelegate struct {
	common    *common.Common
	copiedIdx int
}

// NewItemDelegate creates a new ItemDelegate.
func NewItemDelegate(common *common.Common) *ItemDelegate {
	return &ItemDelegate{
		common:    common,
		copiedIdx: -1,
	}
}

// Width returns the item width.
func (d ItemDelegate) Width() int {
	width := d.common.Styles.MenuItem.GetHorizontalFrameSize() + d.common.Styles.MenuItem.GetWidth()
	return width
}

// Height returns the item height. Implements list.ItemDelegate.
func (d *ItemDelegate) Height() int {
	height := d.common.Styles.MenuItem.GetVerticalFrameSize() + d.common.Styles.MenuItem.GetHeight()
	return height
}

// Spacing returns the spacing between items. Implements list.ItemDelegate.
func (d *ItemDelegate) Spacing() int { return 1 }

// Update implements list.ItemDelegate.
func (d *ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	idx := m.Index()
	item, ok := m.SelectedItem().(gituser.UserRepos)
	if !ok {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.common.KeyMap.Copy):
			d.copiedIdx = idx
			d.common.Output.Copy(item.Command())
			if m.IsFiltered() {
				m.ResetFilter()
				return nil
			}
			return m.SetItem(idx, item)
		}
	}
	return nil
}

func TruncateString(s string, max int) string {
	if max < 0 {
		max = 0
	}
	return truncate.StringWithTail(s, uint(max), "…")
}

// Render implements list.ItemDelegate.
func (d *ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(gituser.UserRepos)
	s := strings.Builder{}
	var matchedRunes []int

	// Conditions
	var (
		isSelected = index == m.Index()
		isFiltered = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	styles := d.common.Styles.RepoSelector.Normal
	if isSelected {
		styles = d.common.Styles.RepoSelector.Active
	}

	title := i.Title()
	title = TruncateString(title, m.Width()-styles.Base.GetHorizontalFrameSize())
	if i.IsPrivate() {
		title += " 🔒"
	}
	if isSelected {
		title += " "
	}
	var updatedStr string
	updatedStr = fmt.Sprintf(" Updated %s", humanize.Time(i.UpdatedAt))

	if m.Width()-styles.Base.GetHorizontalFrameSize()-lipgloss.Width(updatedStr)-lipgloss.Width(title) <= 0 {
		updatedStr = ""
	}
	updatedStyle := styles.Updated.Copy().
		Align(lipgloss.Right).
		Width(m.Width() - styles.Base.GetHorizontalFrameSize() - lipgloss.Width(title))
	updated := updatedStyle.Render(updatedStr)

	if isFiltered && index < len(m.VisibleItems()) {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	if isFiltered {
		unmatched := styles.Title.Copy().Inline(true)
		matched := unmatched.Copy().Underline(true)
		title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
	}
	title = styles.Title.Render(title)
	desc := i.Description()
	desc = TruncateString(desc, m.Width()-styles.Base.GetHorizontalFrameSize())
	desc = styles.Desc.Render(desc)

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, title, updated))
	s.WriteRune('\n')
	s.WriteString(desc)
	s.WriteRune('\n')

	cmd := i.Command()
	cmdStyler := styles.Command.Render
	if d.copiedIdx == index {
		cmd = "(copied to clipboard)"
		cmdStyler = styles.Desc.Render
		// d.copiedIdx = -1
	}
	cmd = TruncateString(cmd, m.Width()-styles.Base.GetHorizontalFrameSize())
	s.WriteString(cmdStyler(cmd))
	fmt.Fprint(w,
		// d.common.Zone.Mark(i.ID(),
		styles.Base.Render(s.String()),
	)
}
