package huhforms

import (
	"time"

	"github.com/charmbracelet/huh"
)

type CreatePasteForm struct {
	content        string
	expirationTime time.Time
}

func NewCreatePasteForm() *huh.Form {
	var form CreatePasteForm
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Content").
				Value(&form.content),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("Paste Expiration").
				Options(huh.NewOptions("1m", "10m", "30m", "1h", "12h", "1d", "1w")...).
				Title("Choose expiraty time").
				Description("After this time the paste will be automaticly deleted"),
		),
	)
}
