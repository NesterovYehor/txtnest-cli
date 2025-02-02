package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds the SSH and backend configuration details
type SSHConfig struct {
	SSHPort    string
	SSHHost    string
	PrivateKey []byte
	BackendURL string
}

// LoadConfig loads the configuration values from environment variables
func LoadSSHConfig() (*SSHConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	// Get environment variables
	sshPort := os.Getenv("SSH_PORT")
	if sshPort == "" {
		sshPort = "2222" // Default SSH Port if not set
	}
	sshHost := os.Getenv("SSH_HOST")
	if sshPort == "" {
		sshPort = "2222" // Default SSH Port if not set
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		return nil, fmt.Errorf("backend URL is required")
	}

	return &SSHConfig{
		SSHHost:    sshHost,
		SSHPort:    sshPort,
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
