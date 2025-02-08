package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Api APIConfig `yaml:"api"`
}

type APIConfig struct {
	BaseUrl    string        `yaml:"base_url"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries int           `yaml:"max_retries"`
}

var appConfig = &Config{} // Initialize to avoid nil pointer

func Init() error {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath != "" {
		viper.AddConfigPath(cfgPath)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Ensure the root folder is checked

	setDefaults()
	bindEnvVars()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("Failed to read config file: %v", err)
		}
		fmt.Println("Config file is not found. Default values used")
	}

	if err := viper.Unmarshal(appConfig); err != nil {
		return fmt.Errorf("Failed to unmarshal config: %v", err)
	}
	appConfig = &Config{
		Api: APIConfig{
			BaseUrl:    viper.GetString("api.base_url"),
			Timeout:    viper.GetDuration("api.timeout"),
			MaxRetries: viper.GetInt("api.max_retries"),
		},
	}

	return nil
}

func Get() *Config {
	return appConfig
}

func setDefaults() {
	viper.SetDefault("api.base_url", "http://localhost:50051/v1")
	viper.SetDefault("api.timeout", 30*time.Second)
	viper.SetDefault("api.max_retries", 5)
}

func bindEnvVars() {
	viper.BindEnv("api.base_url")
	viper.BindEnv("api.timeout")
	viper.BindEnv("api.max_retries")
}
