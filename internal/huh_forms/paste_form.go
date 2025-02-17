package huhforms

import (
	"time"

	"github.com/charmbracelet/huh"
)

type PasteForm struct {
	Title      string
	Content    string
	Expiration time.Duration
}

func NewCreatePasteForm() (*PasteForm, error) {
	var form PasteForm
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Value(&form.Content),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Content").
				Value(&form.Content),
		),
		huh.NewGroup(
			huh.NewSelect[time.Duration]().
				Key("Paste Expiration").
				Options(
					huh.NewOption("1 Hour", time.Hour),
					huh.NewOption("1 Day", 24*time.Hour),
					huh.NewOption("1 Week", 168*time.Hour),
				).
				Title("Choose expiraty time").
				Description("After this time the paste will be automaticly deleted").
				Value(&form.Expiration),
		),
	).Run(); err != nil {
		return nil, err
	}

	return &form, nil
}
