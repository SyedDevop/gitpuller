package progress

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	err      error
	progress progress.Model
	listLen  int
	index    int
	width    int
	height   int
	done     bool
}

var doneStyle = lipgloss.NewStyle().Margin(1, 2)

type (
	DownloadMes string
	ErrMess     struct{ error }
)

func (e ErrMess) Error() string { return e.error.Error() }

func InitialProgress(listLen int) model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	return model{
		listLen:  listLen,
		progress: p,
		index:    0,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case DownloadMes:

		if m.index >= m.listLen-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Quit
		}
		// Update progress bar
		progressCmd := m.progress.SetPercent((float64(m.index) + 1) / float64(m.listLen))
		m.index++
		return m, progressCmd

	case ErrMess:
		m.err = msg
		return m, tea.Quit

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	n := m.listLen
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Downloading %d files/folders.\n", n))
	}
	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)
	prog := m.progress.View()
	return " " + prog + pkgCount
}
