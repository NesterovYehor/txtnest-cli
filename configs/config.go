package config

import (
	"fmt"
	"os"
)

// Config struct holds the SSH and backend configuration details
type SSHConfig struct {
	SSHPort    string
	PrivateKey []byte
	BackendURL string
}

// LoadConfig loads the configuration values from environment variables
func LoadSSHConfig() (*SSHConfig, error) {
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
	sshHostKeyPath := os.Getenv("SSH_HOST_KEY_PATH")
	privateKey, err := os.ReadFile(sshHostKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return &SSHConfig{
		SSHPort:    sshPort,
		PrivateKey: privateKey,
		BackendURL: backendURL,
	}, nil
}

// Config struct holds the SSH and backend configuration details
type CliConfig struct {
	BackendURL string
}

// LoadConfig loads the configuration values from environment variables
func LoadCliConfig() (*CliConfig, error) {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		return nil, fmt.Errorf("backend URL is required")
	}

	return &CliConfig{
		BackendURL: backendURL,
	}, nil
}
