package tui

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PasteEditor allows editing of a paste's content.
type PasteEditor struct {
	textarea textarea.Model
	client   *api.ApiClient
}

// NewPasteEditor returns a new PasteEditor with the given initial content.
func NewPasteEditor(key string, client *api.ApiClient) *PasteEditor {
	paste, err := client.FetchPaste(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ta := textarea.New()
	ta.SetValue(paste.Content)
	ta.Focus()
	ta.CharLimit = 0 // no character limit
	ta.Prompt = "> "
	return &PasteEditor{
		textarea: ta,
		client:   client,
	}
}

// Init returns the initial command for the editor.
func (pe *PasteEditor) Init() tea.Cmd {
	return textarea.Blink
}

// Update updates the textarea.
func (pe *PasteEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	pe.textarea, cmd = pe.textarea.Update(msg)
	return pe, cmd
}

// View renders the editor.
func (pe *PasteEditor) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render("Editing Paste (Ctrl+S to save, Esc to cancel)")
	return fmt.Sprintf("%s\n\n%s", title, pe.textarea.View())
}

// Content returns the current content of the paste.
func (pe *PasteEditor) Content() string {
	return pe.textarea.Value()
}
