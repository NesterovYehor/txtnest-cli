package tui

import (
	"fmt"
	"strings"

	"github.com/NesterovYehor/txtnest-cli/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = itemStyle.Copy().Foreground(lipgloss.Color("12")).Bold(true)
	checkboxStyle     = lipgloss.NewStyle().PaddingRight(1)
)

// PastesList renders a list of paste metadata.
type PastesList struct {
	pastes []models.Metadata
	cursor int
}

// NewPastesList returns a new PastesList.
func NewPastesList(ps []models.Metadata) *PastesList {
	return &PastesList{
		pastes: ps,
		cursor: 0,
	}
}

// Init is part of the Bubbletea Model interface.
func (pl *PastesList) Init() tea.Cmd {
	return nil
}

// Update handles key presses for navigation and selection.
func (pl *PastesList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			if pl.cursor < len(pl.pastes)-1 {
				pl.cursor++
			}
		case "up", "k":
			if pl.cursor > 0 {
				pl.cursor--
			}
		case "enter", " ":
			// Selection confirmed; exit the list view.
			return pl, tea.Quit
		case "q", "ctrl+c":
			return pl, tea.Quit
		}
	}
	return pl, nil
}

// View renders the list of pastes.
func (pl *PastesList) View() string {
	var b strings.Builder

	for i, paste := range pl.pastes {
		checkbox := "[ ]"
		if pl.cursor == i {
			checkbox = "[X]"
		}

		item := fmt.Sprintf("%s %s", checkboxStyle.Render(checkbox), paste.Title)
		if pl.cursor == i {
			item = selectedItemStyle.Render("➤ " + item)
		} else {
			item = itemStyle.Render("  " + item)
		}

		b.WriteString(item)
		if i < len(pl.pastes)-1 {
			b.WriteString("\n")
		}
	}

	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(b.String())
}

// Selected returns the currently highlighted paste.
func (pl *PastesList) Selected() *models.Metadata {
	if len(pl.pastes) == 0 {
		return nil
	}
	return &pl.pastes[pl.cursor]
}

