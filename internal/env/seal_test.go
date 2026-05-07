package env_test

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
)

func sealTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"DB_HOST": "prod-db",
					"API_KEY": "secret",
					"LOG_LEVEL": "warn",
				},
				Sealed: []string{},
			},
		},
	}
}

func TestSeal_SpecificKeys(t *testing.T) {
	cfg := sealTestConfig()
	res, err := env.Seal(cfg, "production", []string{"DB_HOST", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sealed) != 2 {
		t.Errorf("expected 2 sealed keys, got %d", len(res.Sealed))
	}
	set := cfg.EnvSetByName("production")
	if len(set.Sealed) != 2 {
		t.Errorf("expected set.Sealed length 2, got %d", len(set.Sealed))
	}
}

func TestSeal_AllKeys(t *testing.T) {
	cfg := sealTestConfig()
	res, err := env.Seal(cfg, "production", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sealed) != 3 {
		t.Errorf("expected 3 sealed keys, got %d", len(res.Sealed))
	}
}

func TestSeal_DeduplicatesKeys(t *testing.T) {
	cfg := sealTestConfig()
	cfg.EnvSets[0].Sealed = []string{"DB_HOST"}
	res, err := env.Seal(cfg, "production", []string{"DB_HOST", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sealed) != 1 {
		t.Errorf("expected 1 newly sealed key, got %d", len(res.Sealed))
	}
	set := cfg.EnvSetByName("production")
	if len(set.Sealed) != 2 {
		t.Errorf("expected total 2 sealed keys, got %d", len(set.Sealed))
	}
}

func TestSeal_UnknownSet(t *testing.T) {
	cfg := sealTestConfig()
	_, err := env.Seal(cfg, "staging", []string{"DB_HOST"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestSeal_KeyNotInBase(t *testing.T) {
	cfg := sealTestConfig()
	_, err := env.Seal(cfg, "production", []string{"NONEXISTENT"})
	if err == nil {
		t.Fatal("expected error for key not in base")
	}
}

func TestUnseal_SpecificKeys(t *testing.T) {
	cfg := sealTestConfig()
	cfg.EnvSets[0].Sealed = []string{"API_KEY", "DB_HOST", "LOG_LEVEL"}
	res, err := env.Unseal(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sealed) != 1 || res.Sealed[0] != "API_KEY" {
		t.Errorf("expected [API_KEY] unsealed, got %v", res.Sealed)
	}
	set := cfg.EnvSetByName("production")
	if len(set.Sealed) != 2 {
		t.Errorf("expected 2 remaining sealed keys, got %d", len(set.Sealed))
	}
}

func TestUnseal_AllKeys(t *testing.T) {
	cfg := sealTestConfig()
	cfg.EnvSets[0].Sealed = []string{"API_KEY", "DB_HOST"}
	res, err := env.Unseal(cfg, "production", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sealed) != 2 {
		t.Errorf("expected 2 unsealed keys, got %d", len(res.Sealed))
	}
	set := cfg.EnvSetByName("production")
	if len(set.Sealed) != 0 {
		t.Errorf("expected no sealed keys, got %d", len(set.Sealed))
	}
}

func TestUnseal_UnknownSet(t *testing.T) {
	cfg := sealTestConfig()
	_, err := env.Unseal(cfg, "ghost", []string{"DB_HOST"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}
