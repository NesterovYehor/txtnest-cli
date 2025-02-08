package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "txtnest-cli",
	Short: "CLI tool for storing and manage data with encryption and authentication",
	Long:  "A command-line interface for managing pastes securely with encryption and authentication.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to TxtnestCLI! Use --help to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(signupCmd)
}
