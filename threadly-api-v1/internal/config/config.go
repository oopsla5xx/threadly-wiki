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

const (
	LogFormatJSON    = "json"
	LogFormatConsole = "console"

	LogOutputStdout = "stdout"
	LogOutputFile   = "file"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	HTTP      HTTPConfig      `mapstructure:"http"`
	Log       LogConfig       `mapstructure:"log"`
	Databases DatabasesConfig `mapstructure:"databases"`
	Redis     RedisConfig     `mapstructure:"redis"`
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
	Level string 					`mapstructure:"level"`
	Format  string        `mapstructure:"format"`
	Outputs []string      `mapstructure:"outputs"`
	File    LogFileConfig `mapstructure:"file"`
}

type LogFileConfig struct {
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
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

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
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

	errs = append(errs, c.Log.validate()...)

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

	errs = append(errs, c.Redis.validate()...)

	return errors.Join(errs...)
}

func (r *RedisConfig) validate() []error {
	var errs []error

	if r.Addr == "" {
		errs = append(errs, errors.New("redis.addr (REDIS_ADDR) is required"))
	}

	if r.DB < 0 || r.PoolSize < 1 {
		errs = append(errs, fmt.Errorf("redis: need db >= 0 and pool_size >= 1, got db=%d pool_size=%d", r.DB, r.PoolSize))
	}

	return errs
}

func (l *LogConfig) validate() []error {
	var errs []error

	if l.Format != LogFormatJSON && l.Format != LogFormatConsole {
		errs = append(errs, fmt.Errorf("log.format must be %q or %q, got %q", LogFormatJSON, LogFormatConsole, l.Format))
	}

	if len(l.Outputs) == 0 {
		errs = append(errs, errors.New("log.outputs must not be empty"))
	}

	for _, output := range l.Outputs {
		switch output {
		case LogOutputStdout:
		case LogOutputFile:
			if l.File.Path == "" {
				errs = append(errs, errors.New("log.file.path is required when log.outputs contains file"))
			}

			if l.File.MaxSizeMB <= 0 || l.File.MaxBackups < 0 || l.File.MaxAgeDays < 0 {
				errs = append(errs, fmt.Errorf(
					"log.file: need max_size_mb > 0, max_backups >= 0, max_age_days >= 0, got %d/%d/%d",
					l.File.MaxSizeMB, l.File.MaxBackups, l.File.MaxAgeDays,
				))
			}
		default:
			errs = append(errs, fmt.Errorf("log.outputs: unknown output %q", output))
		}
	}

	return errs
}
