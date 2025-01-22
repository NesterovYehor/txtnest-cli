package ui

import (
	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const label = `
 _          _                      _   
| |        | |                    | |  
| |_ __  __| |_  _ __    ___  ___ | |_ 
| __|\ \/ /| __|| '_ \  / _ \/ __|| __|
| |_  >  < | |_ | | | ||  __/\__ \| |_ 
 \__|/_/\_\ \__||_| |_| \___||___/ \__|
    `

const (
	widthCof  = 0.5
	heightCof = 0.6
)

var choices = []list.Item{
	Choice{title: "Create Paste", description: "Create a new paste"},
	Choice{title: "Read Paste", description: "View pastes"},
	Choice{title: "Account", description: "Manage your account"},
	Choice{title: "About", description: "Details about the project"},
	Choice{title: "Exit", description: "Exit the application"},
}

// Choice represents a single menu item.
type Choice struct {
	title       string
	description string
}

// FilterValue implements the list.Item interface.
func (c Choice) FilterValue() string {
	return c.title
}

// Title implements the list.Item interface.
func (c Choice) Title() string {
	return c.title
}

// Description implements the list.Item interface.
func (c Choice) Description() string {
	return c.description
}

// App encapsulates the global state of the application
type App struct {
	client *api.Client
	canvas *Canvas
}

// AppModel represents the main application model
type AppModel struct {
	app        *App
	isMenuMode bool
	menu       list.Model
	content    tea.Model
}

// NewAppModel creates the main application model with the config
func NewAppModel(backendURL string) AppModel {
	cvs := NewCanvas()
	httpClient := api.NewClient(backendURL) // Use the backendURL from config
	menu := list.New(choices, list.NewDefaultDelegate(), 0, 0)
	menu.DisableQuitKeybindings()
	menu.SetShowHelp(false)
	menu.Title = "Menu"
	app := &App{
		canvas: cvs,
		client: httpClient,
	}
	return AppModel{
		menu:       menu,
		app:        app,
		isMenuMode: true,
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
		m.menu.SetSize(width/3, height/2)
		m.app.canvas.Resize(width, height)
		m.content = NewCreatePasteModel(m.app)

	case tea.KeyMsg:
		if !m.isMenuMode && msg.Type == tea.KeyEsc {
			m.isMenuMode = true
			return m, nil
		}
		if m.isMenuMode {
			return m.updateMenu(msg)
		}
		m.content, cmd = m.content.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m AppModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// Handle menu selection dynamically
			switch choice := m.menu.SelectedItem().(Choice); choice.title {
			case "Create Paste":
				m.isMenuMode = false
				return m, nil
			case "Read Paste":
				m.isMenuMode = false
				return m, nil
			case "Account":
				m.isMenuMode = false
				return m, nil
			case "Exit":
				return m, tea.Quit
			}

		case tea.KeyDown, tea.KeyUp:
			m.menu, cmd = m.menu.Update(msg)
			switch choice := m.menu.SelectedItem().(Choice); choice.title {
			case "Create Paste":
				m.content = NewCreatePasteModel(m.app)
				return m, nil
			case "Read Paste":
				m.content = newKeyInputModel(m.app)
				return m, nil
			case "Account":
				m.content = newAccountModel(m.app)
				return m, nil
			case "About":
				m.content = newAboutModel(m.app)
				return m, nil
			}
		}
	}
	return m, cmd
}

// View renders the AppModel's UI
func (m AppModel) View() string {
	if m.app.canvas.width == 0 {
		return "Loading"
	}
	width := int(float32(m.app.canvas.width) / widthCof)
	height := int(float32(m.app.canvas.height) / heightCof)
	menuStyle := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	view := menuStyle.Render(lipgloss.NewStyle().AlignVertical(lipgloss.Top).Render(m.menu.View()))
	view = lipgloss.JoinHorizontal(lipgloss.Center, view, m.content.View())
	view = lipgloss.JoinVertical(lipgloss.Left, label, view)
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		m.app.canvas.Render(view),
	)
}
