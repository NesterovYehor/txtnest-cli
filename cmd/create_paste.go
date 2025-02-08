package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	huhforms "github.com/NesterovYehor/txtnest-cli/internal/huh_forms"
	"github.com/spf13/cobra"
)

var (
	expiration time.Duration
	filePath   string
)

var createCmd = &cobra.Command{
	Use:   "create [content|file]",
	Short: "Create new paste",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.GetInstance()

		// Get content from either:
		// 1. File (-f flag)
		// 2. Direct input (argument)
		// 3. Interactive form (if none provided)
		var contentData []byte
		switch {
		case filePath != "":
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
			contentData = data
		case len(args) > 0:
			contentData = []byte(args[0])
		default:
			// Launch Huh form if no input provided
			form := huhforms.NewCreatePasteForm()
			contentData = []byte(form.Content)
			expiration = form.Expiration
		}

		expTime := time.Now().Add(expiration)
		key, err := client.CreatePaste(expTime, contentData)
		if err != nil {
			fmt.Printf("Error creating paste: %v\n", err)
			return
		}

		fmt.Printf("Paste created! Key: %s\n", key)
	},
}

func init() {
	createCmd.Flags().DurationVarP(&expiration, "expire", "e", 24*time.Hour, "Expiration time")
	createCmd.Flags().StringVarP(&filePath, "file", "f", "", "Read content from file")
}
