package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envctl/envctl/internal/config"
)

const validYAML = `
version: "1"
env_sets:
  - name: base
    variables:
      APP_ENV: production
      LOG_LEVEL: info
  - name: db
    variables:
      DB_HOST: localhost
      DB_PORT: "5432"
targets:
  - name: staging
    envs: [base, db]
  - name: production
    envs: [base, db]
`

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempFile(t, validYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("expected version 1, got %q", cfg.Version)
	}
	if len(cfg.EnvSets) != 2 {
		t.Errorf("expected 2 env sets, got %d", len(cfg.EnvSets))
	}
	if len(cfg.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(cfg.Targets))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/envctl.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_DuplicateEnvSetName(t *testing.T) {
	duplicateYAML := `version: "1"
env_sets:
  - name: base
    variables:
      FOO: bar
  - name: base
    variables:
      BAZ: qux
`
	path := writeTempFile(t, duplicateYAML)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate env set name, got nil")
	}
}

func TestEnvSetByName(t *testing.T) {
	path := writeTempFile(t, validYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	es, err := cfg.EnvSetByName("db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if es.Variables["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", es.Variables["DB_PORT"])
	}

	_, err = cfg.EnvSetByName("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing env set, got nil")
	}
}
