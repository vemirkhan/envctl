package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func pinTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"DB_HOST": "prod.db",
					"API_KEY": "secret",
					"DEBUG":   "false",
				},
				Pinned: []string{},
			},
		},
	}
}

func TestPin_Success(t *testing.T) {
	cfg := pinTestConfig()
	results, err := Pin(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Pinned {
		t.Fatalf("expected key to be pinned, got %+v", results)
	}
	if len(cfg.EnvSets[0].Pinned) != 1 || cfg.EnvSets[0].Pinned[0] != "API_KEY" {
		t.Errorf("expected Pinned=[API_KEY], got %v", cfg.EnvSets[0].Pinned)
	}
}

func TestPin_AlreadyPinned(t *testing.T) {
	cfg := pinTestConfig()
	cfg.EnvSets[0].Pinned = []string{"API_KEY"}
	results, err := Pin(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Errorf("expected skipped for already-pinned key")
	}
}

func TestPin_KeyNotInBase(t *testing.T) {
	cfg := pinTestConfig()
	results, err := Pin(cfg, "production", []string{"MISSING_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped || results[0].Reason != "key not in base" {
		t.Errorf("expected skipped with reason 'key not in base', got %+v", results[0])
	}
}

func TestPin_UnknownSet(t *testing.T) {
	cfg := pinTestConfig()
	_, err := Pin(cfg, "staging", []string{"DB_HOST"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestUnpin_Success(t *testing.T) {
	cfg := pinTestConfig()
	cfg.EnvSets[0].Pinned = []string{"API_KEY", "DEBUG"}
	results, err := Unpin(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Pinned != false || results[0].Skipped {
		t.Errorf("expected unpinned result, got %+v", results[0])
	}
	if len(cfg.EnvSets[0].Pinned) != 1 || cfg.EnvSets[0].Pinned[0] != "DEBUG" {
		t.Errorf("expected Pinned=[DEBUG], got %v", cfg.EnvSets[0].Pinned)
	}
}

func TestUnpin_NotPinned(t *testing.T) {
	cfg := pinTestConfig()
	results, err := Unpin(cfg, "production", []string{"DB_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped || results[0].Reason != "not pinned" {
		t.Errorf("expected skipped with reason 'not pinned', got %+v", results[0])
	}
}
