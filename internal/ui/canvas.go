package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
    border = "#E6EBE0"
    color = "#E6EBE0"
)

type Canvas struct {
	width  int
	height int
	style  lipgloss.Style
}

func NewCanvas() *Canvas {
	return &Canvas{
		style: lipgloss.NewStyle().Bold(true).Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color(border)).Align(lipgloss.Center),
	}
}

func (c *Canvas) Resize(width, height int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	c.width = width
	c.height = height
}

func (c *Canvas) Render(content string) string {
	if c.width == 0 || c.height == 0 {
		return "Loading..."
	}
	return lipgloss.Place(
		c.width,
		c.height,
		lipgloss.Center,
		lipgloss.Center,
		c.style.Width(c.width).Height(c.height).Render(content),
	)
}
