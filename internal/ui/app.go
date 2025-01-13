package ui

import (

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	widthCof  = 0.3
	heightCof = 0.6
)

// App encapsulates the global state of the application
type App struct {
	Client *api.Client
	Canvas *Canvas
}

// AppModel represents the main application model
type AppModel struct {
	app     *App
	content tea.Model
}

// NewApp initializes the App with shared resources
func NewApp() *App {
	return &App{
		Client: api.NewClient("http://52.207.228.123:80"), // Initialize the HTTP client
		Canvas: NewCanvas(),    // Initialize the canvas
	}
}

// NewAppModel creates the main application model
func NewAppModel() AppModel {
    app := NewApp()
	return AppModel{
		app:     app,
		content: NewMenuModle(app),
	}
}

// Init initializes the AppModel
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update updates the AppModel based on incoming messages
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := int(float32(msg.Width) * widthCof)
		height := int(float32(msg.Height) * heightCof)
		m.app.Canvas.Resize(width, height)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.content, cmd = m.content.Update(msg)

	return m, cmd
}

// View renders the AppModel's UI
func (m AppModel) View() string {
	width := int(float32(m.app.Canvas.width) / widthCof)
	height := int(float32(m.app.Canvas.height) / heightCof)
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		m.app.Canvas.Render(m.content.View()),
	)
}
