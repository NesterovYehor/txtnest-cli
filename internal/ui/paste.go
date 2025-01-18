package ui

import (
	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PasteModel struct {
	app      *App
	viewport viewport.Model
	paste    models.Paste
}

func NewPasteModel(app *App, paste models.Paste) *PasteModel {
	vp := viewport.New(int(float32(app.canvas.width)*0.6), int(float32(app.canvas.height)*0.7))
	vp.SetContent(paste.Content)
	vp.GotoTop()
	return &PasteModel{
		app:      app,
		viewport: vp,
		paste:    paste,
	}
}

/* INIT */

func (m PasteModel) Init() tea.Cmd {
	return nil
}

/* VIEW */

func (m PasteModel) View() string {
	// Format the dates as dd/mm/yyyy
	creationDate := m.paste.Creation_date.Format("02/01/2006")
	expirationDate := m.paste.ExpirationDate.Format("02/01/2006")

	// Define styles for the creation and expiration dates
	creationStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(m.app.canvas.width / 4)
	expirationStyle := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.app.canvas.width / 4)

	// Create the top bar with proper spacing between dates
	topBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		creationStyle.Render(creationDate),
		expirationStyle.Render(expirationDate),
	)
	bottomBar := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.app.canvas.width).Render("Esc to close")

	// Style the viewport with borders and padding
	viewportStyle := lipgloss.NewStyle().
		Padding(1, 2).Border(lipgloss.NormalBorder())

	// Combine the top bar and the viewport content
	return lipgloss.JoinVertical(lipgloss.Top, topBar, viewportStyle.Render(m.viewport.View()), bottomBar)
}

/* UPLOAD */

func (m PasteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown, tea.KeyUp:
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
