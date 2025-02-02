package ui

import (
	"fmt"

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

var menu = menuModel{
	Choices: []string{"Create New Paste", "Set Expiration", "Clear Paste", "Exit"},
	Cursor:  0,
}

var expirationMenu = expirationSubMenu{
	Choices: []string{"Never", "1m", "10m", "30m", "1h", "1d", "1w", "1m (1 month)", "1y (1 year)"},
	Cursor:  0,
	Active:  false,
}

type createPasteInput struct {
	textArea       textarea.Model
	expirationMenu expirationSubMenu
}

type createPasteState struct {
	input          *createPasteInput
	inTextAreaMode bool
	selectedExpiry string
	error          error
}

/* VIEW */

// View renders the UI
func (m *model) CreatePasteView() string {
	if m.size == undersized {
		return ""
	} else {
		textAreaView := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(m.state.createPaste.input.textArea.View())
		menuView := m.viewMainMenu()
		expirationView := m.viewExpirationMenu()
		return lipgloss.JoinVertical(lipgloss.Left, "New Paste", textAreaView, menuView, expirationView)
	}
}

/* UPDATE */

// Update handles user input and updates the state
func (m *model) UpdateCreatePaste(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.state.createPaste.inTextAreaMode:
			return m.updateTextArea(msg)
		case expirationMenu.Active:
			return m.updateExpirationMenu(msg)
		default:
			return m.updateMainMenu(msg)
		}
	}
	return m, nil
}

/* VIEW HELP FUNCS*/

func (m *model) viewMainMenu() string {
	if m.state.createPaste.input.expirationMenu.Active {
		return ""
	}

	s := "Menu:\n"
	for i, choice := range menu.Choices {
		cursor := "[ ]"
		if menu.Cursor == i {
			cursor = "[x]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return lipgloss.NewStyle().Render(s)
}

func (m *model) viewExpirationMenu() string {
	if !m.state.createPaste.input.expirationMenu.Active {
		return ""
	}

	s := "Set Expiration:\n"
	for i, choice := range expirationMenu.Choices {
		cursor := "[ ]"
		if expirationMenu.Cursor == i {
			cursor = "[x]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return lipgloss.NewStyle().Render(s)
}

/* UPDATE HELP FUNCS*/

func (m *model) updateTextArea(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyTab:
		m.state.createPaste.inTextAreaMode = false
		m.state.createPaste.input.textArea.Blur()
	default:
		m.state.createPaste.input.textArea, cmd = m.state.createPaste.input.textArea.Update(msg)
	}
	return m, cmd
}

func (m *model) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if menu.Cursor > 0 {
			menu.Cursor--
		}
	case tea.KeyDown:
		if menu.Cursor < len(menu.Choices)-1 {
			menu.Cursor++
		}
	case tea.KeyTab:
		m.state.createPaste.inTextAreaMode = true
		m.state.createPaste.input.textArea.Focus()
		return m, nil
	case tea.KeyEnter:
		switch menu.Cursor {
		case 0: // Save Paste
			m.state.createPaste.input.textArea.Reset()
			return m, nil
		case 1: // Set Expiration
			expirationMenu.Active = true
		case 2: // Clear Text
			m.state.createPaste.input.textArea.Reset()
		case 3:
			m.isMenuActive = true
		}
	}
	return m, nil
}

func (m *model) updateExpirationMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		expirationMenu.Active = false
	case tea.KeyUp:
		if expirationMenu.Cursor > 0 {
			expirationMenu.Cursor--
		}
	case tea.KeyDown:
		if expirationMenu.Cursor < len(expirationMenu.Choices)-1 {
			expirationMenu.Cursor++
		}
	case tea.KeyEnter:
		m.state.createPaste.selectedExpiry = expirationMenu.Choices[expirationMenu.Cursor]
		fmt.Println("Selected expiration:", m.state.createPaste.selectedExpiry)
		expirationMenu.Active = false
	}
	return m, nil
}
