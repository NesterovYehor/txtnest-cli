package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in your txtnext acocunt",
	RunE: func(cmd *cobra.Command, args []string) error {
		var email, password string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email").
					Placeholder("example@email.com").
					Value(&email), // Bind input to email variable
				huh.NewInput().
					Title("Password").
					Placeholder("********").
					EchoMode(huh.EchoModePassword).
					Value(&password), // Bind input to password variable
			),
		)

		// Execute the form
		err := form.Run()
		if err != nil {
			return fmt.Errorf("form failed: %w", err)
		}

		// Call API to create account (uncomment when implemented)
		// err = api.signup(email, username, password)
		// if err != nil {
		// 	return fmt.Errorf("failed to sign up: %w", err)
		// }

		fmt.Println("Log in successful! You can manage your pastes .")
		return nil
	},
}
