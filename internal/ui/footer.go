package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type footerState struct {
	commands []footerCommand
}

type footerCommand struct {
	key   string
	value string
}

func (m model) FooterView() string {
	bold := lipgloss.NewStyle().Bold(true).Render

	table := lipgloss.NewStyle().
		Width(m.widthContainer).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingBottom(1).
		Align(lipgloss.Center)

	commands := []string{}
	for _, cmd := range m.state.footer.commands {
		commands = append(commands, bold(" "+cmd.key+" ")+cmd.value+"  ")
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		table.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				commands...,
			),
		))
}
