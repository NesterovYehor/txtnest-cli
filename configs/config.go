package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds the SSH and backend configuration details
type Config struct {
	SSHPort    string
	PrivateKey []byte
	BackendURL string
}

// LoadConfig loads the configuration values from environment variables
func LoadConfig() (*Config, error) {
	// Load environment variables from .env file located in the config folder
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	// Get environment variables
	sshPort := os.Getenv("SSH_PORT")
	if sshPort == "" {
		sshPort = "2222" // Default SSH Port if not set
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		return nil, fmt.Errorf("backend URL is required")
	}

	// Read private key
	privateKey, err := os.ReadFile("ssh_host_key")
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return &Config{
		SSHPort:    sshPort,
		PrivateKey: privateKey,
		BackendURL: backendURL,
	}, nil
}

