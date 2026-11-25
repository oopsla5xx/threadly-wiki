package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"threadly-api-v1/internal/config"
)

func TestNewWritesJSONToFileWhileStdoutIsConsole(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "app.log")

	appLogger, cleanup, err := New(config.EnvironmentDevelopment, config.LogConfig{
		Level:   "debug",
		Format:  config.LogFormatConsole,
		Outputs: []string{config.LogOutputStdout, config.LogOutputFile},
		File:    config.LogFileConfig{Path: path, MaxSizeMB: 1},
	})
	if err != nil {
		t.Fatal(err)
	}

	appLogger.Debug("hello")
	cleanup()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(content), "\x1b[") {
		t.Fatalf("file contains ANSI color codes: %q", content)
	}

	var entry map[string]any
	if err := json.Unmarshal(content, &entry); err != nil {
		t.Fatalf("file line is not JSON: %v: %q", err, content)
	}

	if entry["msg"] != "hello" || entry["level"] != "debug" || entry["caller"] == nil {
		t.Errorf("unexpected entry: %v", entry)
	}
}

func TestNewFailsFastOnUnwritableFilePath(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(blocker, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	_, _, err := New(config.EnvironmentProduction, config.LogConfig{
		Level:   "info",
		Format:  config.LogFormatJSON,
		Outputs: []string{config.LogOutputFile},
		File:    config.LogFileConfig{Path: filepath.Join(blocker, "app.log"), MaxSizeMB: 1},
	})
	if err == nil {
		t.Fatal("expected error for unwritable log path")
	}
}
