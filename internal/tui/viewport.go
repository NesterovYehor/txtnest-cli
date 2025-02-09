package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	headerHeight = 3
	footerHeight = 2
)

type ContentModel struct {
	content  string
	viewport viewport.Model
	ready    bool
}

func NewContentModel(content string) ContentModel {
	return ContentModel{
		content: content,
	}
}

func (m ContentModel) Init() tea.Cmd {
	return nil
}

func (m ContentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, (msg.Height-headerHeight-footerHeight)/2)
			m.viewport.SetContent(wrapContent(m.content, msg.Width))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ContentModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	header := headerStyle.Render(" Paste Content (q to quit) ")
	body := m.viewport.View()
	footer := footerStyle.Render(fmt.Sprintf(
		" %d%% ", int(m.viewport.ScrollPercent()*100),
	))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	)
}

// Helper to wrap text for viewport
func wrapContent(content string, width int) string {
	return lipgloss.NewStyle().
		Width(width - 2). // account for borders
		Render(content)
}

// Styling
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Bold(true).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Bold(true)
)
