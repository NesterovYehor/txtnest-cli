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
Welcome to TextNest! It's a Pastebin-like app for the terminal, built entirely in Go. The most exciting part of this project for me has been learning how to build the backend from scratch. 

TextNest allows users to create and manage text-based pastes through an anonymous SSH interface. The backend is designed with a microservices architecture, making it modular and scalable. The main goal is to provide fast and reliable access to your pastes while keeping things simple and efficient.

Key Features:
- Anonymous SSH Access: Create and manage pastes without needing an account.
- Microservices Architecture: Different services (API, upload, download) working together seamlessly.
- gRPC & Protobufs: Used for high-performance communication between services.
- Kafka: Helps manage data consistency, especially in cases where several services need to work together or when data safety is a priority.
- Scalable Storage: Blob storage on Amazon S3 and metadata stored in PostgreSQL.
- Redis Caching: Speeds up access to frequently used data.

Technologies:
- Go: The main language used for backend development.
- gRPC & Protobufs: For communication between services.
- Kafka: To ensure reliable data handling and consistency.
- PostgreSQL, Redis, Amazon S3: For data storage and caching.

This is my first backend project, and it’s still a work in progress, with plenty of features yet to be built. I’d love to hear your feedback or advice to help me improve as a developer.

Check out the project repositories:
- Backend repository: https://github.com/NesterovYehor/textnest
- CLI tool repository: https://github.com/NesterovYehor/txtnest-cli

Feel free to reach out if you want to chat or share your thoughts! 
You can find me on Discord: https://discordapp.com/users/591678870973841428 or Twiter(X): https://x.com/_n3st_?s=21.

Thanks for checking out my project!
` // Set the content of the viewport
	view.SetContent(
		lipgloss.NewStyle().
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
