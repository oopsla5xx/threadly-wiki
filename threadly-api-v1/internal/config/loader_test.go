package config

import (
	"strings"
	"testing"
)

func TestLoadMergesBaseOverlayAndEnv(t *testing.T) {
	// configs/ paths are relative to the module root.
	t.Chdir("../..")

	t.Setenv("APP_ENV", EnvironmentProduction)
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("REDIS_ADDR", "redis.internal:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://a.example.com,https://b.example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Databases.Primary.MaxConns != 20 {
		t.Errorf("overlay not merged: max_conns = %d, want 20", cfg.Databases.Primary.MaxConns)
	}
	if cfg.Databases.Primary.MinConns != 2 {
		t.Errorf("base value lost: min_conns = %d, want 2", cfg.Databases.Primary.MinConns)
	}
	if cfg.HTTP.Port != 9090 {
		t.Errorf("env override ignored: port = %d, want 9090", cfg.HTTP.Port)
	}
	if cfg.Databases.Primary.URL != "postgres://u:p@localhost:5432/db" {
		t.Errorf("DATABASE_URL not bound: %q", cfg.Databases.Primary.URL)
	}
	if got := cfg.HTTP.CORS.AllowOrigins; len(got) != 2 || got[0] != "https://a.example.com" || got[1] != "https://b.example.com" {
		t.Errorf("CORS_ALLOW_ORIGINS not split: %q", got)
	}
	if cfg.Redis.Addr != "redis.internal:6380" || cfg.Redis.Password != "secret" {
		t.Errorf("REDIS_ADDR/REDIS_PASSWORD not bound: addr=%q password=%q", cfg.Redis.Addr, cfg.Redis.Password)
	}
	if cfg.Redis.PoolSize != 20 || cfg.Redis.DialTimeout.String() != "5s" {
		t.Errorf("redis overlay/base not merged: %+v", cfg.Redis)
	}
	if cfg.App.Environment != EnvironmentProduction {
		t.Errorf("environment = %q", cfg.App.Environment)
	}
}

func TestLoadDevelopmentLogConfigAndEnvOverrides(t *testing.T) {
	t.Chdir("../..")

	t.Setenv("APP_ENV", EnvironmentDevelopment)
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_OUTPUTS", "stdout,file")
	t.Setenv("LOG_FILE_PATH", "/tmp/threadly/app.log")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	log := cfg.Log
	if log.Format != LogFormatConsole || log.Level != "debug" {
		t.Errorf("dev overlay not merged: format=%q level=%q", log.Format, log.Level)
	}
	if len(log.Outputs) != 2 || log.Outputs[0] != LogOutputStdout || log.Outputs[1] != LogOutputFile {
		t.Errorf("LOG_OUTPUTS not split: %q", log.Outputs)
	}
	if log.File.Path != "/tmp/threadly/app.log" || log.File.MaxSizeMB != 10 || log.File.Compress {
		t.Errorf("file config wrong: %+v", log.File)
	}
	if cfg.Redis.Addr != "localhost:6379" || cfg.Redis.PoolSize != 10 {
		t.Errorf("dev redis defaults wrong: %+v", cfg.Redis)
	}
}

func TestLoadProductionRequiresRedisAddr(t *testing.T) {
	t.Chdir("../..")

	t.Setenv("APP_ENV", EnvironmentProduction)
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://a.example.com")
	t.Setenv("REDIS_ADDR", "")

	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "REDIS_ADDR") {
		t.Fatalf("expected REDIS_ADDR error, got %v", err)
	}
}

func TestValidateRejectsBadLogConfig(t *testing.T) {
	cfg := Config{
		HTTP:      HTTPConfig{Port: 8080, ShutdownTimeout: 1, CORS: CORSConfig{AllowOrigins: []string{"x"}}},
		Databases: DatabasesConfig{Primary: DatabaseConfig{URL: "x", MaxConns: 1}},
		Log:       LogConfig{Format: "xml", Outputs: []string{"file", "syslog"}},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	for _, want := range []string{"log.format", "log.file.path", "max_size_mb", `unknown output "syslog"`} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}
