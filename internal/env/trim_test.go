package env

import (
	"testing"

	"github.com/your-org/envctl/internal/config"
)

func trimTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []*config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"APP_HOST":   "localhost",
					"APP_PORT":   "8080",
					"DB_HOST":    "db",
					"DB_PORT":    "5432",
					"LOG_LEVEL":  "info",
					"LOG_FORMAT": "json",
				},
				Targets: []*config.Target{
					{
						Name: "prod",
						Overrides: map[string]string{
							"APP_HOST": "prod.example.com",
							"DB_HOST":  "prod-db.example.com",
						},
					},
				},
			},
		},
	}
}

func TestTrim_ByPrefix(t *testing.T) {
	cfg := trimTestConfig()
	res, err := Trim(cfg, "app", "APP_", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if _, ok := cfg.EnvSets[0].Base["APP_HOST"]; ok {
		t.Error("APP_HOST should have been removed")
	}
	if _, ok := cfg.EnvSets[0].Base["DB_HOST"]; !ok {
		t.Error("DB_HOST should remain")
	}
}

func TestTrim_BySuffix(t *testing.T) {
	cfg := trimTestConfig()
	res, err := Trim(cfg, "app", "", "_PORT", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if _, ok := cfg.EnvSets[0].Base["APP_PORT"]; ok {
		t.Error("APP_PORT should have been removed")
	}
}

func TestTrim_WithTargets(t *testing.T) {
	cfg := trimTestConfig()
	res, err := Trim(cfg, "app", "APP_", "", []string{"prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 2 from base + 1 from prod target (APP_HOST)
	if len(res.Removed) != 3 {
		t.Errorf("expected 3 removed, got %d: %v", len(res.Removed), res.Removed)
	}
}

func TestTrim_UnknownSet(t *testing.T) {
	cfg := trimTestConfig()
	_, err := Trim(cfg, "missing", "APP_", "", nil)
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestTrim_NoPrefixOrSuffix(t *testing.T) {
	cfg := trimTestConfig()
	_, err := Trim(cfg, "app", "", "", nil)
	if err == nil {
		t.Error("expected error when no prefix or suffix given")
	}
}
