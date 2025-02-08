package cmd

import (
	"fmt"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	"github.com/spf13/cobra"
)

var getPasteCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a paste by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args[0]) != 8 {
			fmt.Printf("Key can't be less or greeater then 8 chars lenth")
			return
		}

		paste, err := api.GetInstance("").FetchPaste(args[0])
		if err != nil {
			fmt.Printf("Error fetching paste: %v\n", err)
			return
		}

		fmt.Printf("Content: %s\nExpires at: %s\n",
			string(paste.Content),
			paste.Metadata.ExpirationDate.Format(time.RFC822))
	},
}
