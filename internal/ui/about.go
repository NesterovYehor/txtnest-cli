package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AboutModel represents the about screen
type aboutModel struct {
	app      *App
	viewport viewport.Model
	message  string
}

// NewAboutModel creates a new AboutModel
func newAboutModel(app *App) aboutModel {
	// Create the viewport with dimensions matching the canvas
	width := int(float32(app.canvas.width) * 0.6)
	height := int(float32(app.canvas.height) * 0.7)
	view := viewport.New(width, height)

	// Define the message for the about screen
	message := `
Welcome to TextNest!

This is a Pastebin-like app for the terminal, fully written in Go. 
The most exciting and important part of this journey for me is learning and building the backend from scratch.

This is my first ever backend and CLI project and the progect isn't finished still and there many things to do,
so I’d love to hear your advice and feedback to help me grow as a developer.

Here are the project links:
- Backend repository: https://github.com/NesterovYehor/textnest
- CLI tool repository: https://github.com/NesterovYehor/txtnest-cli

If you’d like to share your thoughts or connect, feel free to reach out on Discord: Nest.

Thank you for checking out my project!
` // Set the content of the viewport
	view.SetContent(
		lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(height).
			Width(width).
			Render(message),
	)

	return aboutModel{
		app:      app,
		viewport: view,
		message:  message,
	}
}

// Init initializes the AboutModel
func (m aboutModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the AboutModel
func (m aboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the AboutModel's UI
func (m aboutModel) View() string {
	return m.viewport.View()
}
