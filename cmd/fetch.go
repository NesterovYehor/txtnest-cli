package cmd

import (
	"fmt"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/NesterovYehor/txtnest-cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	showViewport bool
	decrypt      bool
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [key]",
	Short: "Retrieve paste by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.GetInstance()
		paste, err := client.FetchPaste(args[0])
		if err != nil {
			fmt.Printf("Error fetching paste: %v\n", err)
			return
		}

		if showViewport {
			// Launch TUI viewport
			model := tui.NewContentModel(string(paste.Content))
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error displaying content: %v\n", err)
			}
		} else {
			// Plain output
			fmt.Printf("Content:\n%s\n", paste.Content)
			fmt.Printf("Expires at: %s\n",
				paste.ExpirationDate.Format(time.RFC822))
		}
	},
}

func init() {
	fetchCmd.Flags().BoolVarP(&showViewport, "view", "v", false, "Show in viewport")
	fetchCmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt content")
}
