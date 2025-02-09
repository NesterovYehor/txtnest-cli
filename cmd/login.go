package cmd

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	huhforms "github.com/NesterovYehor/txtnest-cli/internal/huh_forms"
	"github.com/NesterovYehor/txtnest-cli/internal/storage"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in your txtnext acocunt",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.GetInstance()
		tokenStorage, err := storage.GetTokenStorage()
		if err != nil {
			return err
		}
		// Execute the form
		form, err := huhforms.NewAuthForm()
		if err != nil {
			return fmt.Errorf("form failed: %w", err)
		}

		jwt, err := client.LogIn(form.Email, form.Password)
		if err != nil {
			return fmt.Errorf("failed to log in: %w", err)
		}

		if err := tokenStorage.SaveTokens(jwt); err != nil {
			return err
		}

		fmt.Println("Log in successful! You can manage your pastes.")
		return nil
	},
}
