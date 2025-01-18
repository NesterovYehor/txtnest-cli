package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	renderer *lipgloss.Renderer

	border     lipgloss.TerminalColor
	background lipgloss.TerminalColor
	highlight  lipgloss.TerminalColor
	error      lipgloss.TerminalColor
	body       lipgloss.TerminalColor
	accent     lipgloss.TerminalColor

	base lipgloss.Style
}

func BasicTheme(renderer *lipgloss.Renderer, highlight *string) theme {
	base := theme{
		renderer: renderer,
	}

	base.background = lipgloss.AdaptiveColor{Dark: "#000000", Light: "#FBFCFD"}
	base.border = lipgloss.AdaptiveColor{Dark: "#3A3F42", Light: "#D7DBDF"}
	base.body = lipgloss.AdaptiveColor{Dark: "#889096", Light: "#889096"}
	base.accent = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}
	if highlight != nil {
		base.highlight = lipgloss.Color(*highlight)
	} else {
		base.highlight = lipgloss.Color("#FF5C00")
	}
	base.error = lipgloss.Color("203")

	base.base = renderer.NewStyle().Foreground(base.body)

	return base
}
