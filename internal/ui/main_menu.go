package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

// MenuModel manages the state of the menu.
type menuState struct {
	list list.Model
}

// View renders the menu.
func (m model) MenuView() string {
	return m.state.menu.list.View()
}

// Update handles user input.
func (m model) UpdateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyDown, tea.KeyUp:
			// First, update the list so the selection changes
			m.state.menu.list, cmd = m.state.menu.list.Update(msg)

			// Now fetch the updated selection
			selectedItem, ok := m.state.menu.list.SelectedItem().(Choice)
			if ok {
				switch selectedItem.title {
				case "Create Paste":
					m = m.SwitchPage(createPage)
				case "Read Paste":
					m = m.SwitchPage(inputKeyPage)
				case "Account":
					m = m.SwitchPage(accountPage)
				case "About":
					m = m.SwitchPage(aboutPage)
				case "Exit":
					return m, tea.Quit
				}
			}

		default:
			// Process other keys normally
			m.state.menu.list, cmd = m.state.menu.list.Update(msg)
		}
	}

	return m, cmd
}
