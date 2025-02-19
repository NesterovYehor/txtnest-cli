package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/NesterovYehor/txtnest-cli/config"
	"github.com/NesterovYehor/txtnest-cli/internal/models"
	"github.com/NesterovYehor/txtnest-cli/internal/validation"
)

type TokenUpdateCallback func(newTokens *models.TokenData) error

var (
	instance *ApiClient
	once     sync.Once
)

type ApiClient struct {
	baseUrl         string
	httpClient      *http.Client
	accessToken     string
	refreshToken    string
	tokenExp        time.Time
	tokenMu         sync.Mutex
	updateCallbacks []TokenUpdateCallback
}

func (c *ApiClient) RegisterTokenUpdateCallback(cb TokenUpdateCallback) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.updateCallbacks = append(c.updateCallbacks, cb)
}

func (c *ApiClient) SetTokens(tokens *models.TokenData) error {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// First clear existing tokens
	c.accessToken = ""
	c.refreshToken = ""
	c.tokenExp = time.Time{}

	// Validate token hierarchy
	if tokens == nil || tokens.RefreshToken == "" {
		return errors.New("no valid tokens provided")
	}

	// Check refresh token validity first
	if !validation.ValidateRefreshToken(tokens.RefreshExpiresAt) {
		return errors.New("refresh token expired")
	}

	// Check if access token needs refresh
	if !validation.ValidateAccessToken(tokens.ExpiresAt) {
		if err := c.refreshTokens(); err != nil {
			return fmt.Errorf("failed to refresh tokens: %w", err)
		}
		return nil
	}

	// If we get here, use the provided valid tokens
	c.accessToken = tokens.AccessToken
	c.refreshToken = tokens.RefreshToken
	c.tokenExp = tokens.ExpiresAt
	return nil
}

func GetInstance() *ApiClient {
	apiCfg := config.Get().Api
	once.Do(func() {
		instance = &ApiClient{
			baseUrl: apiCfg.BaseUrl,
			httpClient: &http.Client{
				Timeout: apiCfg.Timeout,
				Transport: &http.Transport{
					MaxIdleConns:       100,
					IdleConnTimeout:    90 * time.Second,
					DisableCompression: true,
				},
			},
		}
	})
	return instance
}

func (c *ApiClient) CreatePaste(title string, expirationDate time.Time, content []byte) (string, error) {
	start := time.Now()
	newPaste := map[string]any{
		"expiration_date": expirationDate,
		"title":           title,
	}
	ress, err := c.makeRequest("POST", "/upload", newPaste, nil)
	if err != nil {
		return "", err
	}
	defer ress.Body.Close()
	var output struct {
		Key       string `json:"key"`
		UploadURL string `json:"upload_url"`
	}
	if err := json.NewDecoder(ress.Body).Decode(&output); err != nil {
		return "", fmt.Errorf("Failed to decode http response: %v", err)
	}
	buf := bytes.NewBufferString(string(content))
	_, err = c.makeContentRequest("PUT", output.UploadURL, buf)
	if err != nil {
		return "", fmt.Errorf("Failed to create new paste: %v", err)
	}
	// Calculate and print the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Request and decode took: %v\n", elapsed)
	return output.Key, nil
}

func (c *ApiClient) FetchPaste(key string) (*models.Paste, error) {
	start := time.Now()

	ress, err := c.makeRequest("GET", "/download", map[string]any{"key": key}, nil)
	if err != nil {
		return nil, err
	}
	defer ress.Body.Close()

	var paste models.Paste
	if err := json.NewDecoder(ress.Body).Decode(&paste); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %v", err)
	}

	contentRess, err := c.makeContentRequest("GET", paste.ContentURL, nil)
	if err != nil {
		return nil, err
	}
	defer contentRess.Body.Close()

	content, err := io.ReadAll(contentRess.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %v", err)
	}
	paste.Content = string(content)

	fmt.Printf("FetchPaste took: %v\n", time.Since(start))

	return &paste, nil
}

func (c *ApiClient) SignUp(email, name, password string) error {
	user := map[string]any{
		"name":     name,
		"email":    email,
		"password": password,
	}
	_, err := c.makeRequest("POST", "/signup", user, nil)
	if err != nil {
		return fmt.Errorf("Failed registration: %v", err)
	}
	return nil
}

func (c *ApiClient) LogIn(email, password string) (*models.TokenData, error) {
	user := map[string]any{
		"email":    email,
		"password": password,
	}
	resp, err := c.makeRequest("GET", "/login", user, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed registration: %v", err)
	}
	var jwt models.TokenData
	if err := json.NewDecoder(resp.Body).Decode(&jwt); err != nil {
		return nil, fmt.Errorf("Failed to decode response:%v", err)
	}
	return &jwt, nil
}

func (c *ApiClient) FetchAllTokens() ([]models.Metadata, error) {
	resp, err := c.makeRequest("GET", "/download/all", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch list of pastes: %w", err)
	}
	defer resp.Body.Close()
	var output struct {
		Pastes []models.Metadata `json:"pastes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}
	fmt.Println(output)

	return output.Pastes, nil
}

func (c *ApiClient) UpdatePaste(key, content string) error {
	output, err := c.makeRequest("GET", fmt.Sprintf("/update/%v", key), nil, nil)
	if err != nil {
		return err
	}
	var response struct {
		Url string `json:"update_url"`
	}
	if err := json.NewDecoder(output.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding JSON response: %w", err)
	}
	buf := bytes.NewBufferString(string(content))
	fmt.Println(buf)
	_, err = c.makeContentRequest("PUT", response.Url, buf)
	if err != nil {
		return fmt.Errorf("Failed to create new paste: %v", err)
	}
	return nil
}

func (c *ApiClient) refreshTokens() error {
	resp, err := c.makeRequest("POST", "/refresh", map[string]string{
		"refresh_token": c.refreshToken,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to refresh tokens: %v", err)
	}
	defer resp.Body.Close()

	var newTokens models.TokenData
	if err := json.NewDecoder(resp.Body).Decode(&newTokens); err != nil {
		return fmt.Errorf("Failed to decode response:%v", err)
	}

	for _, cd := range c.updateCallbacks {
		if err := cd(&newTokens); err != nil {
			fmt.Printf("token update callback failed: %v", err)
		}
	}
	return nil
}

func (c *ApiClient) makeRequest(method, endpoint string, body any, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseUrl+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new http request: %v", err)
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to process http request: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(errBody))
	}

	return resp, nil
}

func (c *ApiClient) makeContentRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(errBody))
	}

	return resp, nil
}
