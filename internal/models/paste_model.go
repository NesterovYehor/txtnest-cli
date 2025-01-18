package models

import "time"

type Paste struct {
	Content        string    `json:"content"`
	Creation_date  time.Time `json:"creation_date"`
	ExpirationDate time.Time `json:"expiration_date"`
}
