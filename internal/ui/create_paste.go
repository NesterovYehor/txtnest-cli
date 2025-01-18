package ui

import (
	"fmt"
	"sync"

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
	expirationMenu expirationSubMenu
	inTextAreaMode bool
	selectedExpiry string
	error          error
}

// Singleton instance
var (
	instance *CreatePasteModel
	once     sync.Once
)

// NewCreatePasteModel initializes the CreatePasteModel as a singleton
func NewCreatePasteModel(app *App) CreatePasteModel {
	once.Do(func() {
		ta := textarea.New()
		ta.Placeholder = "Enter text..."
		ta.ShowLineNumbers = false
		ta.SetWidth(int(float32(app.canvas.width) * 0.6))
		ta.SetHeight(int(float32(app.canvas.height) * 0.6))
		ta.Prompt = ""
		ta.CharLimit = 1000000
		ta.Focus()

		instance = &CreatePasteModel{
			textArea: ta,
			app:      app,
			menu: menuModel{
				Choices: []string{"Create New Paste", "Set Expiration", "Clear Paste", "Exit"},
				Cursor:  0,
			},
			expirationMenu: expirationSubMenu{
				Choices: []string{"Never", "1m", "10m", "30m", "1h", "1d", "1w", "1m (1 month)", "1y (1 year)"},
				Cursor:  0,
				Active:  false,
			},
			inTextAreaMode: true,
		}
	})
	return *instance
} /* INIT */

// Init initializes the model
func (m CreatePasteModel) Init() tea.Cmd {
	return m.textArea.Cursor.BlinkCmd()
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
	return lipgloss.NewStyle().Render(s)
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
	return lipgloss.NewStyle().Render(s)
}

/* UPDATE HELP FUNCS*/

func (m CreatePasteModel) updateTextArea(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyCtrlS:
		m.inTextAreaMode = false
		m.textArea.Blur()
	default:
		m.textArea, cmd = m.textArea.Update(msg)
	}
	return m, cmd
}

func (m CreatePasteModel) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.menu.Cursor > 0 {
			m.menu.Cursor--
		}
	case tea.KeyDown:
		if m.menu.Cursor < len(m.menu.Choices)-1 {
			m.menu.Cursor++
		}
	case tea.KeyTab:
		m.inTextAreaMode = true
		m.textArea.Focus()
		return m, nil
	case tea.KeyEnter:
		switch m.menu.Cursor {
		case 0: // Save Paste
			key, err := m.app.client.CreatePaste(m.textArea.Value(), m.expirationMenu.Choices[m.expirationMenu.Cursor])
			if err != nil {
				m.error = err
				fmt.Println(err)
			}
			m.textArea.Reset()
			fmt.Println(key)
			return m, nil
		case 1: // Set Expiration
			m.expirationMenu.Active = true
		case 2: // Clear Text
			m.textArea.Reset()
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
