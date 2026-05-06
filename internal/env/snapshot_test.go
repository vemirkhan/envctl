package env_test

import (
	"testing"
	"time"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
)

func snapshotTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{"PORT": "8080", "DEBUG": "false"},
				Targets: map[string]map[string]string{
					"prod": {"DEBUG": "false", "PORT": "443"},
				},
			},
		},
		Snapshots: []config.SnapshotEntry{},
	}
}

func TestTakeSnapshot_Basic(t *testing.T) {
	cfg := snapshotTestConfig()
	snap, err := env.TakeSnapshot(cfg, "app", "", "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Name != "snap1" {
		t.Errorf("expected name snap1, got %s", snap.Name)
	}
	if snap.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", snap.Vars["PORT"])
	}
	if len(cfg.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot in config, got %d", len(cfg.Snapshots))
	}
}

func TestTakeSnapshot_WithTarget(t *testing.T) {
	cfg := snapshotTestConfig()
	snap, err := env.TakeSnapshot(cfg, "app", "prod", "prod-snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Vars["PORT"] != "443" {
		t.Errorf("expected PORT=443, got %s", snap.Vars["PORT"])
	}
	if snap.Target != "prod" {
		t.Errorf("expected target prod, got %s", snap.Target)
	}
}

func TestTakeSnapshot_DuplicateName(t *testing.T) {
	cfg := snapshotTestConfig()
	_, _ = env.TakeSnapshot(cfg, "app", "", "dup")
	_, err := env.TakeSnapshot(cfg, "app", "", "dup")
	if err == nil {
		t.Fatal("expected error for duplicate snapshot name")
	}
}

func TestTakeSnapshot_UnknownSet(t *testing.T) {
	cfg := snapshotTestConfig()
	_, err := env.TakeSnapshot(cfg, "missing", "", "s")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestListSnapshots_All(t *testing.T) {
	cfg := snapshotTestConfig()
	cfg.Snapshots = []config.SnapshotEntry{
		{Name: "a", Set: "app", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{Name: "b", Set: "app", CreatedAt: time.Now().Add(-1 * time.Hour)},
	}
	list := env.ListSnapshots(cfg, "")
	if len(list) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(list))
	}
	if list[0].Name != "a" {
		t.Errorf("expected sorted by time, first should be a")
	}
}

func TestDeleteSnapshot_Success(t *testing.T) {
	cfg := snapshotTestConfig()
	cfg.Snapshots = []config.SnapshotEntry{{Name: "to-del", Set: "app"}}
	if err := env.DeleteSnapshot(cfg, "to-del"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Snapshots) != 0 {
		t.Errorf("expected 0 snapshots after delete")
	}
}

func TestDeleteSnapshot_NotFound(t *testing.T) {
	cfg := snapshotTestConfig()
	if err := env.DeleteSnapshot(cfg, "ghost"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
