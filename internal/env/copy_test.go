package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func copyTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "base",
				Base: map[string]string{
					"APP_ENV": "development",
					"LOG_LEVEL": "debug",
				},
			},
			{
				Name: "staging",
				Base: map[string]string{
					"APP_ENV": "staging",
				},
				Targets: map[string]map[string]string{
					"us-east": {"REGION": "us-east-1"},
				},
			},
			{
				Name: "production",
				Base: map[string]string{},
			},
		},
	}
}

func TestCopy_BaseOnly(t *testing.T) {
	cfg := copyTestConfig()
	result, err := Copy(cfg, "base", "production", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SourceSet != "base" || result.DestSet != "production" {
		t.Errorf("unexpected result sets: %+v", result)
	}
	if len(result.Keys) != 2 {
		t.Errorf("expected 2 keys copied, got %d", len(result.Keys))
	}
}

func TestCopy_WithTarget(t *testing.T) {
	cfg := copyTestConfig()
	result, err := Copy(cfg, "staging", "production", "us-east")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// staging base has APP_ENV, target adds REGION => 2 keys
	if len(result.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(result.Keys), result.Keys)
	}
}

func TestCopy_UnknownSource(t *testing.T) {
	cfg := copyTestConfig()
	_, err := Copy(cfg, "nope", "production", "")
	if err == nil {
		t.Fatal("expected error for unknown source")
	}
}

func TestCopy_UnknownDestination(t *testing.T) {
	cfg := copyTestConfig()
	_, err := Copy(cfg, "base", "nope", "")
	if err == nil {
		t.Fatal("expected error for unknown destination")
	}
}

func TestCopy_DestinationReceivesValues(t *testing.T) {
	cfg := copyTestConfig()
	_, err := Copy(cfg, "base", "production", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod, _ := cfg.EnvSetByName("production")
	if prod.Base["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug in destination, got %q", prod.Base["LOG_LEVEL"])
	}
}
