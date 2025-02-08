package huhforms

import (
	"github.com/charmbracelet/huh"
)

type RegistrationForm struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewRegistrationForm() (*RegistrationForm, error) {
	var form RegistrationForm
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Placeholder("example@email.com").
				Value(&form.Email), // Bind input to email variable
			huh.NewInput().
				Title("Username").
				Placeholder("unique user name").
				Value(&form.Name), // Bind input to username variable
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
	return &form, err
}
