package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func reorderTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"HOST": "localhost",
					"PORT": "8080",
					"DEBUG": "false",
				},
				KeyOrder: []string{"HOST", "PORT", "DEBUG"},
			},
		},
	}
}

func TestReorder_Success(t *testing.T) {
	cfg := reorderTestConfig()
	err := Reorder(cfg, "app", []string{"DEBUG", "HOST", "PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("app")
	if set.KeyOrder[0] != "DEBUG" || set.KeyOrder[1] != "HOST" || set.KeyOrder[2] != "PORT" {
		t.Errorf("unexpected key order: %v", set.KeyOrder)
	}
}

func TestReorder_PartialOrder(t *testing.T) {
	cfg := reorderTestConfig()
	err := Reorder(cfg, "app", []string{"PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("app")
	if set.KeyOrder[0] != "PORT" {
		t.Errorf("expected PORT first, got %v", set.KeyOrder)
	}
	if len(set.KeyOrder) != 3 {
		t.Errorf("expected 3 keys total, got %d", len(set.KeyOrder))
	}
}

func TestReorder_UnknownSet(t *testing.T) {
	cfg := reorderTestConfig()
	err := Reorder(cfg, "nope", []string{"HOST"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestReorder_UnknownKey(t *testing.T) {
	cfg := reorderTestConfig()
	err := Reorder(cfg, "app", []string{"MISSING"})
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestReorder_DuplicateKey(t *testing.T) {
	cfg := reorderTestConfig()
	err := Reorder(cfg, "app", []string{"HOST", "HOST"})
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
}
