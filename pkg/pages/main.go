package pages

import "github.com/charmbracelet/lipgloss"

type Page interface {
	Title() string
	Render() string
}

type Model struct {
	theme       Theme
	renderer    *lipgloss.Renderer
	pages       []Page
	currentPage int
	width       int
	height      int
}

func NewPageModel() *Model {
	renderer := lipgloss.DefaultRenderer()

	return &Model{
		renderer:    renderer,
		currentPage: 0,
		theme:       GetTheme(renderer),

		pages: []Page{
			&RepoPage{},
		},
	}
}
