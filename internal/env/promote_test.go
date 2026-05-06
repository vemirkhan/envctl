package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func promoteTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{"HOST": "localhost", "PORT": "8080"},
				Targets: map[string]map[string]string{
					"staging": {"HOST": "staging.example.com", "DEBUG": "true"},
					"production": {"HOST": "prod.example.com"},
				},
			},
		},
	}
}

func TestPromote_Success(t *testing.T) {
	cfg := promoteTestConfig()
	res, err := Promote(cfg, "app", "staging", "production", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FromTarget != "staging" || res.ToTarget != "production" {
		t.Errorf("unexpected targets in result: %+v", res)
	}
	// DEBUG should be promoted; HOST should not (no overwrite)
	prod := cfg.EnvSets[0].Targets["production"]
	if prod["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true in production, got %q", prod["DEBUG"])
	}
	if prod["HOST"] != "prod.example.com" {
		t.Errorf("HOST should not be overwritten, got %q", prod["HOST"])
	}
}

func TestPromote_Overwrite(t *testing.T) {
	cfg := promoteTestConfig()
	_, err := Promote(cfg, "app", "staging", "production", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := cfg.EnvSets[0].Targets["production"]
	if prod["HOST"] != "staging.example.com" {
		t.Errorf("expected HOST overwritten to staging value, got %q", prod["HOST"])
	}
}

func TestPromote_UnknownSet(t *testing.T) {
	cfg := promoteTestConfig()
	_, err := Promote(cfg, "nope", "staging", "production", false)
	if err == nil {
		t.Fatal("expected error for unknown env set")
	}
}

func TestPromote_UnknownFromTarget(t *testing.T) {
	cfg := promoteTestConfig()
	_, err := Promote(cfg, "app", "dev", "production", false)
	if err == nil {
		t.Fatal("expected error for unknown source target")
	}
}

func TestPromote_UnknownToTarget(t *testing.T) {
	cfg := promoteTestConfig()
	_, err := Promote(cfg, "app", "staging", "canary", false)
	if err == nil {
		t.Fatal("expected error for unknown destination target")
	}
}

func TestPromote_KeysPromotedList(t *testing.T) {
	cfg := promoteTestConfig()
	res, err := Promote(cfg, "app", "staging", "production", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only DEBUG should be in promoted list (HOST skipped, no overwrite)
	if len(res.KeysPromoted) != 1 || res.KeysPromoted[0] != "DEBUG" {
		t.Errorf("expected [DEBUG] promoted, got %v", res.KeysPromoted)
	}
}
