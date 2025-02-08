package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Register a new Pastebin account",
	RunE: func(cmd *cobra.Command, args []string) error {
		var email, username, password string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email").
					Placeholder("example@email.com").
					Value(&email), // Bind input to email variable
				huh.NewInput().
					Title("Username").
					Placeholder("unique user name").
					Value(&username), // Bind input to username variable
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

		fmt.Println("Signup successful! You can now login.")
		return nil
	},
}
