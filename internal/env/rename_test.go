package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func renameTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "base",
				Vars: map[string]string{"FOO": "bar", "BAZ": "qux"},
			},
			{
				Name: "staging",
				Vars: map[string]string{"FOO": "staging-bar"},
				Targets: []config.Target{{Ref: "base", Path: ".env"}},
			},
		},
	}
}

func TestRename_Success(t *testing.T) {
	cfg := renameTestConfig()
	result, err := Rename(cfg, "base", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OldName != "base" || result.NewName != "production" {
		t.Errorf("unexpected result names: %+v", result)
	}
	if result.KeysUpdated != 2 {
		t.Errorf("expected 2 keys updated, got %d", result.KeysUpdated)
	}
	if cfg.EnvSets[0].Name != "production" {
		t.Errorf("expected env set name to be 'production', got %q", cfg.EnvSets[0].Name)
	}
}

func TestRename_UpdatesTargetRefs(t *testing.T) {
	cfg := renameTestConfig()
	_, err := Rename(cfg, "base", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvSets[1].Targets[0].Ref != "production" {
		t.Errorf("expected target ref to be updated to 'production', got %q", cfg.EnvSets[1].Targets[0].Ref)
	}
}

func TestRename_NotFound(t *testing.T) {
	cfg := renameTestConfig()
	_, err := Rename(cfg, "nonexistent", "newname")
	if err == nil {
		t.Fatal("expected error for unknown env set")
	}
}

func TestRename_DestinationExists(t *testing.T) {
	cfg := renameTestConfig()
	_, err := Rename(cfg, "base", "staging")
	if err == nil {
		t.Fatal("expected error when destination name already exists")
	}
}

func TestRename_SameName(t *testing.T) {
	cfg := renameTestConfig()
	_, err := Rename(cfg, "base", "base")
	if err == nil {
		t.Fatal("expected error when old and new names are the same")
	}
}
