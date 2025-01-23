package main

import (
	"fmt"
	"log/slog"
	"os"

	config "github.com/NesterovYehor/txtnest-cli/configs"
	"github.com/NesterovYehor/txtnest-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	log, err := os.Create("output.log")
	if err != nil {
		panic(err)
	}
	defer log.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(log, &slog.HandlerOptions{})))

	cfg, err := config.LoadCliConfig()
	if err != nil {
		fmt.Printf("log: %v\n", err)
		return
	}

	model := ui.NewAppModel(cfg.BackendURL)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
