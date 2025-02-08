package cmd

import (
	"fmt"

	"github.com/NesterovYehor/txtnest-cli/internal/api"
	huhforms "github.com/NesterovYehor/txtnest-cli/internal/huh_forms"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Register a new Pastebin account",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.GetInstance()

		// Execute the form
		form, err := huhforms.NewRegistrationForm()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = client.SignUp(form.Email, form.Name, form.Password)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Signup successful! You can now login.")
	},
}
