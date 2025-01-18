package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NesterovYehor/txtnest-cli/internal/models"
)

type Paste struct {
	Content    string
	Expiration time.Time
}

func (c *Client) CreatePaste(content string, expiration string) (string, error) {
	body := map[string]any{
		"content":         content,
		"expiration_date": calculateExpiration(expiration),
	}

	resp, err := c.doRequest("POST", "/v1/upload", body, nil)
	if err != nil {
		return "", err
	}

	var response struct {
		Key string `json:"key"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("Failed to decode htpt response, err: %v", err)
	}

	return response.Key, nil
}

func (c *Client) FetchPaste(key string) (*models.Paste, error) {
	body := map[string]any{
		"key": key,
	}

	resp, err := c.doRequest("GET", "/v1/download", body, nil)
	if err != nil {
		return nil, err
	}

	var paste models.Paste

	err = json.NewDecoder(resp.Body).Decode(&paste)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode htpt response, err: %v", err)
	}
	return &paste, err
}

func calculateExpiration(choice string) time.Time {
	switch choice {
	case "Never":
		return time.Time{}
	case "1m":
		return time.Now().Add(1 * time.Minute)
	case "10m":
		return time.Now().Add(10 * time.Minute)
	case "30m":
		return time.Now().Add(30 * time.Minute)
	case "1h":
		return time.Now().Add(1 * time.Hour)
	case "1d":
		return time.Now().Add(24 * time.Hour)
	case "1w":
		return time.Now().Add(7 * 24 * time.Hour)
	case "1m (1 month)":
		return time.Now().AddDate(0, 1, 0) // Add 1 month
	case "1y (1 year)":
		return time.Now().AddDate(1, 0, 0) // Add 1 year
	default:
		// Return zero time if the choice is invalid
		return time.Time{}
	}
}
