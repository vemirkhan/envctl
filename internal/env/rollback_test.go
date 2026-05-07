package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func rollbackTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{"APP_ENV": "prod", "LOG_LEVEL": "warn"},
				Targets: []config.Target{
					{Name: "us-east", Overrides: map[string]string{"REGION": "us-east-1"}},
				},
			},
		},
		Snapshots: []config.Snapshot{
			{
				Name:   "v1",
				EnvSet: "production",
				Vars:   map[string]string{"APP_ENV": "staging", "LOG_LEVEL": "debug"},
				Targets: map[string]map[string]string{
					"us-east": {"REGION": "us-east-2"},
				},
			},
		},
	}
}

func TestRollback_RestoresBase(t *testing.T) {
	cfg := rollbackTestConfig()
	res, err := Rollback(cfg, "production", "v1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.EnvSet != "production" || res.SnapshotName != "v1" {
		t.Errorf("unexpected result metadata: %+v", res)
	}
	set := cfg.EnvSetByName("production")
	if set.Base["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", set.Base["APP_ENV"])
	}
	if set.Base["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", set.Base["LOG_LEVEL"])
	}
}

func TestRollback_RestoresTarget(t *testing.T) {
	cfg := rollbackTestConfig()
	res, err := Rollback(cfg, "production", "v1", "us-east")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Restored["REGION"] != "us-east-2" {
		t.Errorf("expected REGION=us-east-2, got %q", res.Restored["REGION"])
	}
	set := cfg.EnvSetByName("production")
	if set.Targets[0].Overrides["REGION"] != "us-east-2" {
		t.Errorf("target override not restored")
	}
}

func TestRollback_UnknownSet(t *testing.T) {
	cfg := rollbackTestConfig()
	_, err := Rollback(cfg, "nope", "v1", "")
	if err == nil {
		t.Error("expected error for unknown env set")
	}
}

func TestRollback_UnknownSnapshot(t *testing.T) {
	cfg := rollbackTestConfig()
	_, err := Rollback(cfg, "production", "ghost", "")
	if err == nil {
		t.Error("expected error for unknown snapshot")
	}
}

func TestRollback_UnknownTarget(t *testing.T) {
	cfg := rollbackTestConfig()
	_, err := Rollback(cfg, "production", "v1", "eu-west")
	if err == nil {
		t.Error("expected error for unknown target in snapshot")
	}
}
