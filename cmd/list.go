package cmd

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/NesterovYehor/txtnest-cli/internal/storage"
	"github.com/NesterovYehor/txtnest-cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// ListCmd is the Cobra command for listing and editing pastes.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Retrieve list of your pastes",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.GetInstance()

		tokensStorage, err := storage.GetTokenStorage()
		if err != nil {
			fmt.Println("Error getting token storage:", err)
			return
		}
		if tokens, err := tokensStorage.GetTokens(); err == nil {
			if err := client.SetTokens(tokens); err != nil {
				fmt.Printf("Failed to set tokens: %v\n", err)
				return
			}
		}

		pastes, err := client.FetchAllTokens()
		if err != nil {
			fmt.Println("Error fetching pastes:", err)
			return
		}

		model := tui.NewPasteBrowser(client, pastes)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error displaying list of your pastes: %v\n", err)
			return
		}
	},
}

