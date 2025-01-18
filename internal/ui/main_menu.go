package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// MenuModel manages the state of the menu.
type MenuModel struct {
	menu list.Model
}

// Init initializes the MenuModel.
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// View renders the menu.
func (m *MenuModel) View() string {
	return m.menu.View()
}

// Update handles user input.
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.menu.Index() {
			case 0:
				// Implement NewCreatePasteModel logic
			case 1:
				// Implement NewCodeInputModel logic
			case 2:
				// Handle account logic
			case 3:
				return m, tea.Quit
			}
		default:
			m.menu, cmd = m.menu.Update(msg)
			return m, nil
		}
	}
	return m, cmd
}
