package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func cloneTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"APP_ENV": "production",
					"LOG_LEVEL": "warn",
				},
				Targets: []config.Target{
					{
						Name:      "us-east",
						Overrides: map[string]string{"REGION": "us-east-1"},
					},
				},
			},
			{
				Name: "staging",
				Base: map[string]string{
					"APP_ENV": "staging",
				},
			},
		},
	}
}

func TestClone_BaseOnly(t *testing.T) {
	cfg := cloneTestConfig()
	res, err := Clone(cfg, "production", "production-copy", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysCopied != 2 {
		t.Errorf("expected 2 keys copied, got %d", res.KeysCopied)
	}
	dst := cfg.EnvSetByName("production-copy")
	if dst == nil {
		t.Fatal("destination env set not found")
	}
	if len(dst.Targets) != 0 {
		t.Errorf("expected no targets, got %d", len(dst.Targets))
	}
	if dst.Base["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", dst.Base["APP_ENV"])
	}
}

func TestClone_WithTargets(t *testing.T) {
	cfg := cloneTestConfig()
	res, err := Clone(cfg, "production", "prod-backup", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysCopied != 2 {
		t.Errorf("expected 2 keys copied, got %d", res.KeysCopied)
	}
	dst := cfg.EnvSetByName("prod-backup")
	if dst == nil {
		t.Fatal("destination env set not found")
	}
	if len(dst.Targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(dst.Targets))
	}
	if dst.Targets[0].Overrides["REGION"] != "us-east-1" {
		t.Errorf("expected REGION=us-east-1, got %s", dst.Targets[0].Overrides["REGION"])
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	cfg := cloneTestConfig()
	_, err := Clone(cfg, "nonexistent", "copy", false)
	if err == nil {
		t.Fatal("expected error for unknown source")
	}
}

func TestClone_DestinationExists(t *testing.T) {
	cfg := cloneTestConfig()
	_, err := Clone(cfg, "production", "staging", false)
	if err == nil {
		t.Fatal("expected error when destination already exists")
	}
}

func TestClone_IsolatesBase(t *testing.T) {
	cfg := cloneTestConfig()
	_, err := Clone(cfg, "production", "prod-iso", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mutate original; clone should not be affected
	src := cfg.EnvSetByName("production")
	src.Base["APP_ENV"] = "mutated"
	dst := cfg.EnvSetByName("prod-iso")
	if dst.Base["APP_ENV"] != "production" {
		t.Errorf("clone was affected by mutation of source: got %s", dst.Base["APP_ENV"])
	}
}
