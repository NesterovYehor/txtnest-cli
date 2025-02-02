package ui

import (
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyInputState struct {
	codeInput textinput.Model
}

func(m model) newKeyInputModel() keyInputState {
	ti := textinput.New()
	ti.Placeholder = "Enter a code of a paste"
	ti.CharLimit = 8
	ti.Width = 12
	ti.Focus()
	inputKeyState := keyInputState{
		codeInput: ti,
	}
	return inputKeyState
}

/* VIEW */
func (m model) InputKeyView() string {
	// Center the TextInput within the viewport
	codeInputView := lipgloss.Place(
		m.widthContent,
		m.heightContent,
		lipgloss.Center,
		lipgloss.Center,
		m.state.keyInput.codeInput.View(),
	)
	return codeInputView
}

/* UPDATE */
func (m *model) InputKeyUpadte(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state.keyInput.codeInput.Reset()
			m.state.readPaste = *m.newReadPaste(models.Paste{
				Content:        "hkhjfklajsdfkadsjfajfklajfklajsdflkajsdfkljasfkljakls",
				Creation_date:  time.Now(),
				ExpirationDate: time.Now(),
			})
			m.page = readPage
			return m, nil
		default:
			m.state.keyInput.codeInput, cmd = m.state.keyInput.codeInput.Update(msg)
		}
	}
	return m, cmd
}
