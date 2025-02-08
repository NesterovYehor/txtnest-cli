package huhforms

import "github.com/charmbracelet/huh"

type AuthenticationForm struct {
	Email    string
	Password string
}

func NewAuthForm() (*AuthenticationForm, error) {
	var form AuthenticationForm
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Placeholder("example@email.com").
				Value(&form.Email), // Bind input to email variable
			huh.NewInput().
				Title("Password").
				Placeholder("********").
				EchoMode(huh.EchoModePassword).
				Value(&form.Password), // Bind input to password variable
		),
	).Run()
	if err != nil {
		return nil, err
	}
	return &form, nil
}
