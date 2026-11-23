package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env is optional
	}

	environment := os.Getenv("APP_ENV")

	if environment == "" {
		return nil, fmt.Errorf("APP_ENV is required")
	}

	if environment != "development" && environment != "production" {
		return nil, fmt.Errorf("invalid APP_ENV: %s", environment)
	}

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = fmt.Sprintf("configs/config.%s.yaml", environment)
	}

	configLoader := viper.New()

	configLoader.SetConfigFile(configFile)
	configLoader.SetConfigType("yaml")
	configLoader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configLoader.AutomaticEnv()

	environmentBindings := map[string]string{
		"app.name":                          "APP_NAME",
		"app.environment":                   "APP_ENV",
		"http.host":                         "HTTP_HOST",
		"http.port":                         "HTTP_PORT",
		"http.read_timeout":                 "HTTP_READ_TIMEOUT",
		"http.write_timeout":                "HTTP_WRITE_TIMEOUT",
		"http.idle_timeout":                 "HTTP_IDLE_TIMEOUT",
		"databases.primary.url":             "DATABASE_PRIMARY_URL",
		"databases.primary.max_conns":       "DATABASE_PRIMARY_MAX_CONNS",
		"databases.primary.min_conns":       "DATABASE_PRIMARY_MIN_CONNS",
		"databases.primary.max_conn_lifetime": "DATABASE_PRIMARY_MAX_CONN_LIFETIME",
		"databases.primary.max_conn_idle_time": "DATABASE_PRIMARY_MAX_CONN_IDLE_TIME",
		"databases.analytics.url":           "DATABASE_ANALYTICS_URL",
		"databases.analytics.max_conns":      "DATABASE_ANALYTICS_MAX_CONNS",
		"databases.analytics.min_conns":      "DATABASE_ANALYTICS_MIN_CONNS",
		"databases.analytics.max_conn_lifetime": "DATABASE_ANALYTICS_MAX_CONN_LIFETIME",
		"databases.analytics.max_conn_idle_time": "DATABASE_ANALYTICS_MAX_CONN_IDLE_TIME",
		"redis.host":                        "REDIS_HOST",
		"redis.port":                        "REDIS_PORT",
		"redis.database":                    "REDIS_DATABASE",
	}

	for configKey, environmentKey := range environmentBindings {
		if err := configLoader.BindEnv(configKey, environmentKey); err != nil {
			return nil, fmt.Errorf("bind %s: %w", environmentKey, err)
		}
	}

	if err := configLoader.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config

	if err := configLoader.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}
