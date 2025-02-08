package models

import "time"

type Paste struct {
	ExpirationDate time.Time `json:"expiration_date"`
	CreationDate   time.Time `json:"creation_date"`
	Content        string    `json:"content"`
}
