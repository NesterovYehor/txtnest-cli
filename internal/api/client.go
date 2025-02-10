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
		// Use existing refresh token to get new tokens
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

func (c *ApiClient) CreatePaste(expirationDate time.Time, content []byte) (string, error) {
	start := time.Now()
	newPaste := map[string]any{
		"expiration_date": expirationDate,
		"content":         content,
	}
	ress, err := c.makeReuest("POST", "/upload", newPaste, nil)
	if err != nil {
		return "", err
	}
	defer ress.Body.Close()
	var output struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(ress.Body).Decode(&output); err != nil {
		return "", fmt.Errorf("Failed to decode http response: %v", err)
	}
	// Calculate and print the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Request and decode took: %v\n", elapsed)
	return output.Key, nil
}

func (c *ApiClient) FetchPaste(key string) (*models.Paste, error) {
	// Record the start time
	start := time.Now()

	// Make the request
	ress, err := c.makeReuest("GET", "/download", map[string]any{"key": key}, nil)
	if err != nil {
		return nil, err
	}
	defer ress.Body.Close()

	// Decode the response into a Paste model
	var paste *models.Paste
	if err := json.NewDecoder(ress.Body).Decode(&paste); err != nil {
		fmt.Println(ress)
		return nil, fmt.Errorf("failed to decode http response: %v", err)
	}

	// Calculate and print the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Request and decode took: %v\n", elapsed)

	return paste, nil
}

func (c *ApiClient) SignUp(email, name, password string) error {
	user := map[string]any{
		"name":     name,
		"email":    email,
		"password": password,
	}
	_, err := c.makeReuest("POST", "/signup", user, nil)
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
	resp, err := c.makeReuest("GET", "/login", user, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed registration: %v", err)
	}
	var jwt models.TokenData
	if err := json.NewDecoder(resp.Body).Decode(&jwt); err != nil {
		return nil, fmt.Errorf("Failed to decode response:%v", err)
	}
	return &jwt, nil
}

func (c *ApiClient) refreshTokens() error {
	c.tokenMu.Lock()

	defer c.tokenMu.Unlock()

	resp, err := c.makeReuest("POST", "/refresh", map[string]string{
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

func (c *ApiClient) makeReuest(method, endpoint string, body any, headers map[string]string) (*http.Response, error) {
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
		fmt.Println(c.baseUrl + endpoint)
		return nil, fmt.Errorf("Failed to process http request: %v, Status Code:%v", err, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		errBody, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(errBody))
	}

	return resp, nil
}
