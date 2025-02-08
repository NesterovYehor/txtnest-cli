package models

import "time"

type Paste struct {
	Metadata Metadata
	Content  []byte
}

type Metadata struct {
	Key            string    `json:"key"`
	ExpirationDate time.Time `json:"expiration_date"`
	CreationDate   time.Time `json:"creation_date"`
}
