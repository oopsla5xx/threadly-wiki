package config

import (
	"os"
	"testing"
)

func TestLoadMergesBaseOverlayAndEnv(t *testing.T) {
	// configs/ paths are relative to the module root.
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}

	t.Setenv("APP_ENV", EnvironmentProduction)
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("HTTP_PORT", "9090")
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
	if cfg.App.Environment != EnvironmentProduction {
		t.Errorf("environment = %q", cfg.App.Environment)
	}
}
