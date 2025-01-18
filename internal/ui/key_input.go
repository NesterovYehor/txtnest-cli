package ui

import (
	"sync"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keyInputInstance *keyInputModel
	keyInputOnce     sync.Once
)

type keyInputModel struct {
	app       *App
	codeInput textinput.Model
}

func newKeyInputModel(app *App) *keyInputModel {
	keyInputOnce.Do(func() {
		ti := textinput.New()
		ti.Placeholder = "Enter a code of a paste"
		ti.CharLimit = 8
		ti.Width = 12
		ti.Focus()
		keyInputInstance = &keyInputModel{
			app:       app,
			codeInput: ti,
		}
	})
	return keyInputInstance
}

/* INIT */
func (m *keyInputModel) Init() tea.Cmd {
	return nil
}

/* VIEW */
func (m *keyInputModel) View() string {
	codeInputView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		PaddingLeft(2).
		PaddingRight(2).
		MarginBottom(1).
		Render(m.codeInput.View())

	menu := lipgloss.NewStyle().PaddingRight(5).Render("Esc to return")
	if len([]rune(m.codeInput.Value())) == 8 {
		menu = lipgloss.JoinHorizontal(lipgloss.Center, menu, "Enter to read paste")
	}

	return lipgloss.JoinVertical(lipgloss.Center, codeInputView, menu)
}

/* UPDATE */
func (m *keyInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			paste, err := m.app.client.FetchPaste(m.codeInput.Value())
			if err != nil {
				panic(err)
			}
			m.codeInput.Reset()
			return NewPasteModel(m.app, *paste), nil
		}
	}
	m.codeInput, cmd = m.codeInput.Update(msg)
	return m, cmd
}

