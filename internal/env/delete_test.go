package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func deleteTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "base",
				Vars: map[string]string{"APP_ENV": "production"},
			},
			{
				Name: "staging",
				Vars: map[string]string{"APP_ENV": "staging"},
			},
		},
		Targets: []config.Target{
			{
				Name: "deploy",
				Overrides: []config.Override{
					{EnvSet: "staging", Vars: map[string]string{"EXTRA": "1"}},
				},
			},
		},
	}
}

func TestDelete_Success(t *testing.T) {
	cfg := deleteTestConfig()
	if err := Delete(cfg, "staging", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := cfg.EnvSetByName("staging"); err == nil {
		t.Error("expected staging to be deleted")
	}
	if len(cfg.EnvSets) != 1 {
		t.Errorf("expected 1 env set, got %d", len(cfg.EnvSets))
	}
}

func TestDelete_NotFound(t *testing.T) {
	cfg := deleteTestConfig()
	err := Delete(cfg, "nonexistent", false)
	if err == nil {
		t.Fatal("expected error for missing env set")
	}
}

func TestDelete_RemovesRefs(t *testing.T) {
	cfg := deleteTestConfig()
	if err := Delete(cfg, "staging", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, tgt := range cfg.Targets {
		for _, ov := range tgt.Overrides {
			if ov.EnvSet == "staging" {
				t.Errorf("found dangling override referencing deleted set in target %q", tgt.Name)
			}
		}
	}
}

func TestDelete_PreservesOtherSets(t *testing.T) {
	cfg := deleteTestConfig()
	if err := Delete(cfg, "staging", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := cfg.EnvSetByName("base"); err != nil {
		t.Error("expected base env set to still exist")
	}
}

func TestDelete_KeepsRefsWhenFlagFalse(t *testing.T) {
	cfg := deleteTestConfig()
	if err := Delete(cfg, "staging", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Override should still be present since removeRefs=false.
	if len(cfg.Targets[0].Overrides) != 1 {
		t.Errorf("expected overrides to be untouched, got %d", len(cfg.Targets[0].Overrides))
	}
}
