package ui

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuModel struct {
	Choices []string
	Cursor  int
}

type expirationSubMenu struct {
	Choices []string
	Cursor  int
	Active  bool
}

type CreatePasteModel struct {
	textArea       textarea.Model
	app            *App
	menu           menuModel
	httpClient     *api.Client
	expirationMenu expirationSubMenu
	inTextAreaMode bool
	selectedExpiry string
	error          error
}

// NewCreatePasteModel initializes the CreatePasteModel
func NewCreatePasteModel(app *App) CreatePasteModel {
	ta := textarea.New()
	ta.Placeholder = "Enter text..."
	ta.ShowLineNumbers = false
	ta.SetWidth(int(float32(app.Canvas.width) * 0.9))
	ta.SetHeight(int(float32(app.Canvas.height) * 0.9))
	ta.Prompt = ""
	ta.Focus()

	return CreatePasteModel{
		textArea:   ta,
		app:        app,
		menu: menuModel{
			Choices: []string{"Create New Paste", "Clear Paste", "Set Expiration", "Exit"},
			Cursor:  0,
		},
		expirationMenu: expirationSubMenu{
			Choices: []string{"Never", "1m", "10m", "30m", "1h", "1d", "1w", "1m (1 month)", "1y (1 year)"},
			Cursor:  0,
			Active:  false,
		},
		inTextAreaMode: true,
	}
}

/* INIT */

// Init initializes the model
func (m CreatePasteModel) Init() tea.Cmd {
	return nil
}

/* VIEW */

// View renders the UI
func (m CreatePasteModel) View() string {
	textAreaView := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(m.textArea.View())
	menuView := m.viewMainMenu()
	expirationView := m.viewExpirationMenu()
	return lipgloss.JoinVertical(lipgloss.Left, "New Paste", textAreaView, menuView, expirationView)
}

/* UPDATE */

// Update handles user input and updates the state
func (m CreatePasteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.inTextAreaMode:
			return m.updateTextArea(msg)
		case m.expirationMenu.Active:
			return m.updateExpirationMenu(msg)
		default:
			return m.updateMainMenu(msg)
		}
	}
	return m, nil
}

/* VIEW HELP FUNCS*/

func (m CreatePasteModel) viewMainMenu() string {
	if m.expirationMenu.Active {
		return ""
	}

	s := "Menu:\n"
	for i, choice := range m.menu.Choices {
		cursor := "[ ]"
		if m.menu.Cursor == i {
			cursor = "[x]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(s)
}

func (m CreatePasteModel) viewExpirationMenu() string {
	if !m.expirationMenu.Active {
		return ""
	}

	s := "Set Expiration:\n"
	for i, choice := range m.expirationMenu.Choices {
		cursor := "[ ]"
		if m.expirationMenu.Cursor == i {
			cursor = "[x]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(s)
}

/* UPDATE HELP FUNCS*/

func (m CreatePasteModel) updateTextArea(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyEsc:
		m.inTextAreaMode = false
		m.textArea.Blur()
	default:
		m.textArea, cmd = m.textArea.Update(msg)
	}
	return m, cmd
}

func (m CreatePasteModel) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyUp:
		if m.menu.Cursor > 0 {
			m.menu.Cursor--
		}
	case tea.KeyDown:
		if m.menu.Cursor < len(m.menu.Choices)-1 {
			m.menu.Cursor++
		}
	case tea.KeyEnter:
		switch m.menu.Cursor {
		case 0: // Save Paste
			key, err := m.httpClient.CreatePaste(m.textArea.Value(), m.expirationMenu.Choices[m.expirationMenu.Cursor])
			if err != nil {
				m.error = err
			}
			fmt.Println("Paste saved with content:", m.textArea.Value(), "Expiration:", m.selectedExpiry, "The kay to get a paste is:", key)
		case 1: // Clear Text
			m.textArea.Reset()
		case 2: // Set Expiration
			m.expirationMenu.Active = true
		case 3: // Exit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m CreatePasteModel) updateExpirationMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.expirationMenu.Active = false
	case tea.KeyUp:
		if m.expirationMenu.Cursor > 0 {
			m.expirationMenu.Cursor--
		}
	case tea.KeyDown:
		if m.expirationMenu.Cursor < len(m.expirationMenu.Choices)-1 {
			m.expirationMenu.Cursor++
		}
	case tea.KeyEnter:
		m.selectedExpiry = m.expirationMenu.Choices[m.expirationMenu.Cursor]
		fmt.Println("Selected expiration:", m.selectedExpiry)
		m.expirationMenu.Active = false
	}
	return m, nil
}
