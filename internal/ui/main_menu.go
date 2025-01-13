package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	menuChoices []string
	app         *App
	cursor      int
}

func NewMenuModle(app *App) MenuModel {
	return MenuModel{
		menuChoices: []string{"Create Paste", "Read Paste", "Exit"},
		app:      app,
		cursor:      0,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) View() string {
	s := "Main Menu:\n\n"
	for i, choice := range m.menuChoices {
		cursor := " "
		if m.cursor == i {
			cursor = "->"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.menuChoices)-1 {
				m.cursor++
			}
		case "esc":
			return nil, tea.Quit
		case "enter":
			switch m.cursor {
			case 0:
				return NewCreatePasteModel(m.app), nil
			case 1:
				fmt.Println(1)
			case 2:
				fmt.Println(2)
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
