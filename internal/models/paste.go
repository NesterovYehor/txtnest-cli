package models

import (
	"time"
)

type Paste struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Title          string    `json:"title"`
	CreationDate   time.Time `json:"creation_date"`
	ContentURL     string    `json:"content_url"`
	Content        string    `json:"content"`
}

type Metadata struct {
	Key            string    `json:"key"`
	Title          string    `json:"title"`
	CreationDate   time.Time `json:"creation_date"`
	ExpirationDate time.Time `json:"expiration_date"`
}
