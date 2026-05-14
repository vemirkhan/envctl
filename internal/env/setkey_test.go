package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func setkeyTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{"EXISTING": "yes"},
				Targets: []config.Target{
					{Name: "prod", Overrides: map[string]string{"PROD_ONLY": "1"}},
				},
			},
		},
	}
}

func TestSetKey_NewBaseKey(t *testing.T) {
	cfg := setkeyTestConfig()
	if err := SetKey(cfg, "app", "NEW_KEY", "hello", SetKeyOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvSets[0].Base["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", cfg.EnvSets[0].Base["NEW_KEY"])
	}
}

func TestSetKey_OverwriteBase(t *testing.T) {
	cfg := setkeyTestConfig()
	if err := SetKey(cfg, "app", "EXISTING", "new", SetKeyOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvSets[0].Base["EXISTING"] != "new" {
		t.Errorf("expected EXISTING=new")
	}
}

func TestSetKey_NoOverwriteConflict(t *testing.T) {
	cfg := setkeyTestConfig()
	err := SetKey(cfg, "app", "EXISTING", "bad", SetKeyOptions{})
	if err == nil {
		t.Fatal("expected error for duplicate key without overwrite")
	}
}

func TestSetKey_TargetOverride(t *testing.T) {
	cfg := setkeyTestConfig()
	if err := SetKey(cfg, "app", "NEW_PROD", "42", SetKeyOptions{Target: "prod"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvSets[0].Targets[0].Overrides["NEW_PROD"] != "42" {
		t.Errorf("expected NEW_PROD=42 in prod target")
	}
}

func TestSetKey_UnknownSet(t *testing.T) {
	cfg := setkeyTestConfig()
	err := SetKey(cfg, "ghost", "K", "v", SetKeyOptions{})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestSetKey_UnknownTarget(t *testing.T) {
	cfg := setkeyTestConfig()
	err := SetKey(cfg, "app", "K", "v", SetKeyOptions{Target: "staging"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestSetKey_EmptyKey(t *testing.T) {
	cfg := setkeyTestConfig()
	err := SetKey(cfg, "app", "", "v", SetKeyOptions{})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}
