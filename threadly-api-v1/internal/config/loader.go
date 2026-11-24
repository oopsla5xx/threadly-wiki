package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const baseConfigFile = "configs/config.yaml"

// environmentBindings is the complete list of environment variables that
// override config values. Keep it in sync with .env.example.
var environmentBindings = map[string]string{
	"databases.primary.url": "DATABASE_URL",
	"http.host":             "HTTP_HOST",
	"http.port":             "HTTP_PORT",
	"log.level":             "LOG_LEVEL",
}

// Load reads configs/config.yaml, merges the environment overlay
// (configs/config.<APP_ENV>.yaml, or CONFIG_FILE if set) on top,
// then applies environment variable overrides.
func Load() (*Config, error) {
	// .env is optional; real environment variables take precedence.
	_ = godotenv.Load()

	environment := os.Getenv("APP_ENV")

	if environment == "" {
		return nil, fmt.Errorf("APP_ENV is required")
	}

	if environment != EnvironmentDevelopment && environment != EnvironmentProduction {
		return nil, fmt.Errorf("invalid APP_ENV: %s", environment)
	}

	overlayFile := os.Getenv("CONFIG_FILE")

	if overlayFile == "" {
		overlayFile = fmt.Sprintf("configs/config.%s.yaml", environment)
	}

	configLoader := viper.New()
	configLoader.SetConfigType("yaml")

	configLoader.SetConfigFile(baseConfigFile)

	if err := configLoader.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read %s: %w", baseConfigFile, err)
	}

	configLoader.SetConfigFile(overlayFile)

	if err := configLoader.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("merge %s: %w", overlayFile, err)
	}

	for configKey, environmentKey := range environmentBindings {
		if err := configLoader.BindEnv(configKey, environmentKey); err != nil {
			return nil, fmt.Errorf("bind %s: %w", environmentKey, err)
		}
	}

	var config Config

	if err := configLoader.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	config.App.Environment = environment

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
