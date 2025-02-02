package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/NesterovYehor/txtnest-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	log, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer log.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(log, &slog.HandlerOptions{})))

	model, err := ui.NewModel(lipgloss.DefaultRenderer())
	if err != nil {
		panic(err)
	}
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
