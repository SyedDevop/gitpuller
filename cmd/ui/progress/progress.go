package progress

import (
	"fmt"

	"github.com/SyedDevop/gitpuller/cmd/api"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	err      error
	packages []api.Repo
	spinner  spinner.Model
	progress progress.Model
	index    int
	width    int
	height   int
	done     bool
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

type (
	DownloadMes string
	ErrMess     struct{ error }
)

func (e ErrMess) Error() string { return e.error.Error() }

func InitialProgress(list []api.Repo) model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return model{
		packages: list,
		spinner:  s,
		progress: p,
		index:    0,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Printf("%s %s", checkMark, m.packages[m.index].Name),
		m.spinner.Tick,
	)
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
		if m.index >= len(m.packages)-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Quit
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent((float64(m.index) + 1) / float64(len(m.packages)))

		m.index++

		batch := tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.packages[m.index].Name), // print success message above our program
		)
		return m, batch

	case ErrMess:
		m.err = msg
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
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
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Downloading %d files/folders.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index+1, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.packages[m.index].Name)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Downloading " + pkgName)

	// cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	// gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + " " + prog + pkgCount
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
