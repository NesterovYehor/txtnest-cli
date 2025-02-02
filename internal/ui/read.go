package ui

import (
	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type readPasteState struct {
	viewport viewport.Model
	paste    models.Paste
}

func (m model) newReadPaste(paste models.Paste) *readPasteState {
	vp := viewport.New(m.widthContent, m.heightContent)
	vp.SetContent(paste.Content)
	vp.GotoTop()
	return &readPasteState{
		viewport: vp,
		paste:    paste,
	}
}

func (m model) ReadPasteView() string {
	// Format the dates as dd/mm/yyyy
	creationDate := m.state.readPaste.paste.Creation_date.Format("02/01/2006")
	expirationDate := m.state.readPaste.paste.ExpirationDate.Format("02/01/2006")

	// Define styles for the creation and expiration dates
	creationStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(4)
	expirationStyle := lipgloss.NewStyle().Align(lipgloss.Right).Width(4)

	// Create the top bar with proper spacing between dates
	topBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		creationStyle.Render(creationDate),
		expirationStyle.Render(expirationDate),
	)

	// Style the viewport with borders and padding
	viewportStyle := lipgloss.NewStyle().
		Padding(1, 2).Border(lipgloss.NormalBorder())

	// Combine the top bar and the viewport content
	return lipgloss.JoinVertical(lipgloss.Top, topBar, viewportStyle.Render(m.viewport.View()))
}

func (m model) ReadPasteUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.page = inputKeyPage
			return m, nil
		case tea.KeyDown, tea.KeyUp:
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
