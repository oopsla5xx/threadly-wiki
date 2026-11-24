package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	HTTP      HTTPConfig      `mapstructure:"http"`
	Log       LogConfig       `mapstructure:"log"`
	Databases DatabasesConfig `mapstructure:"databases"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	// Environment is set from APP_ENV, not from the config files.
	Environment string `mapstructure:"-"`
}

type HTTPConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	CORS            CORSConfig    `mapstructure:"cors"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type DatabasesConfig struct {
	Primary DatabaseConfig `mapstructure:"primary"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

// Validate reports every invalid field at once so startup fails with a complete message.
func (c *Config) Validate() error {
	var errs []error

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		errs = append(errs, fmt.Errorf("http.port must be in 1-65535, got %d", c.HTTP.Port))
	}

	if c.HTTP.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("http.shutdown_timeout must be positive"))
	}

	if len(c.HTTP.CORS.AllowOrigins) == 0 {
		errs = append(errs, errors.New("CORS_ALLOW_ORIGINS is required"))
	}

	primary := c.Databases.Primary

	if primary.URL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	if primary.MinConns < 0 || primary.MaxConns < 1 || primary.MinConns > primary.MaxConns {
		errs = append(errs, fmt.Errorf(
			"databases.primary: need 0 <= min_conns <= max_conns and max_conns >= 1, got min=%d max=%d",
			primary.MinConns, primary.MaxConns,
		))
	}

	return errors.Join(errs...)
}
