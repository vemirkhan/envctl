package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func protectTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"DB_URL": "postgres://prod",
					"API_KEY": "secret",
					"LOG_LEVEL": "warn",
				},
				Protected: []string{},
			},
		},
	}
}

func TestProtect_Success(t *testing.T) {
	cfg := protectTestConfig()
	err := Protect(cfg, "production", []string{"DB_URL", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("production")
	if len(set.Protected) != 2 {
		t.Errorf("expected 2 protected keys, got %d", len(set.Protected))
	}
}

func TestProtect_Deduplicates(t *testing.T) {
	cfg := protectTestConfig()
	_ = Protect(cfg, "production", []string{"DB_URL"})
	_ = Protect(cfg, "production", []string{"DB_URL", "API_KEY"})
	set := cfg.EnvSetByName("production")
	if len(set.Protected) != 2 {
		t.Errorf("expected 2 protected keys, got %d", len(set.Protected))
	}
}

func TestProtect_UnknownSet(t *testing.T) {
	cfg := protectTestConfig()
	err := Protect(cfg, "staging", []string{"DB_URL"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestProtect_UnknownKey(t *testing.T) {
	cfg := protectTestConfig()
	err := Protect(cfg, "production", []string{"MISSING_KEY"})
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestUnprotect_Success(t *testing.T) {
	cfg := protectTestConfig()
	_ = Protect(cfg, "production", []string{"DB_URL", "API_KEY"})
	err := Unprotect(cfg, "production", []string{"DB_URL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("production")
	if len(set.Protected) != 1 || set.Protected[0] != "API_KEY" {
		t.Errorf("expected only API_KEY protected, got %v", set.Protected)
	}
}

func TestUnprotect_UnknownSet(t *testing.T) {
	cfg := protectTestConfig()
	err := Unprotect(cfg, "staging", []string{"DB_URL"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestProtectedKeys_ReturnsList(t *testing.T) {
	cfg := protectTestConfig()
	_ = Protect(cfg, "production", []string{"LOG_LEVEL"})
	keys, err := ProtectedKeys(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 || keys[0] != "LOG_LEVEL" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
