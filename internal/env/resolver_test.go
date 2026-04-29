package env

import (
	"os"
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func baseConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "backend",
				Vars: map[string]string{
					"APP_ENV":   "development",
					"LOG_LEVEL": "debug",
				},
				Targets: map[string]map[string]string{
					"production": {
						"APP_ENV":   "production",
						"LOG_LEVEL": "warn",
					},
				},
			},
		},
	}
}

func TestResolve_BaseOnly(t *testing.T) {
	cfg := baseConfig()
	resolved, err := Resolve(cfg, "backend", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved["APP_ENV"] != "development" {
		t.Errorf("expected APP_ENV=development, got %q", resolved["APP_ENV"])
	}
	if resolved["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", resolved["LOG_LEVEL"])
	}
}

func TestResolve_WithTarget(t *testing.T) {
	cfg := baseConfig()
	resolved, err := Resolve(cfg, "backend", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", resolved["APP_ENV"])
	}
	if resolved["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL=warn, got %q", resolved["LOG_LEVEL"])
	}
}

func TestResolve_UnknownSet(t *testing.T) {
	cfg := baseConfig()
	_, err := Resolve(cfg, "nonexistent", "")
	if err == nil {
		t.Fatal("expected error for unknown env set, got nil")
	}
}

func TestResolve_UnknownTarget(t *testing.T) {
	cfg := baseConfig()
	_, err := Resolve(cfg, "backend", "staging")
	if err == nil {
		t.Fatal("expected error for unknown target, got nil")
	}
}

func TestResolve_EnvExpansion(t *testing.T) {
	os.Setenv("TEST_SECRET", "supersecret")
	t.Cleanup(func() { os.Unsetenv("TEST_SECRET") })

	cfg := &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Vars: map[string]string{
					"DB_PASS": "${TEST_SECRET}",
				},
			},
		},
	}
	resolved, err := Resolve(cfg, "app", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved["DB_PASS"] != "supersecret" {
		t.Errorf("expected DB_PASS=supersecret, got %q", resolved["DB_PASS"])
	}
}

func TestToExportLines(t *testing.T) {
	r := ResolvedEnv{"FOO": "bar", "BAZ": "it's alive"}
	lines := r.ToExportLines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestResolve_NilConfig(t *testing.T) {
	_, err := Resolve(nil, "backend", "")
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}
