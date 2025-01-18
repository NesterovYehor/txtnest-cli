package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AccountModel represents the account screen
type accountModel struct {
	app      *App
	viewport viewport.Model
	message  string
}

// NewAccountModel creates a new AccountModel
func newAccountModel(app *App) accountModel {
	// Create the viewport with dimensions matching the canvas
	width := int(float32(app.canvas.width) * 0.6)
	height := int(float32(app.canvas.height) * 0.7)
	view := viewport.New(width, height)

	// Define the message for the account screen
	message := "Account and paste management is in development stage."

	// Set the content of the viewport
	view.SetContent(
		lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(height).
			Width(width).
			Render(message),
	)

	return accountModel{
		app:      app,
		viewport: view,
		message:  message,
	}
}

// Init initializes the AccountModel
func (m accountModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the AccountModel
func (m accountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the AccountModel's UI
func (m accountModel) View() string {
	return m.viewport.View()
}

