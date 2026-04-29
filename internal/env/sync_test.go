package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envctl/internal/config"
)

func syncTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Vars: map[string]string{
					"BASE_URL": "https://example.com",
					"DEBUG":    "false",
				},
				Targets: []config.Target{
					{Name: "staging", File: "staging.env", Format: "dotenv",
						Overrides: map[string]string{"DEBUG": "true"}},
					{Name: "prod", File: "prod.env", Format: "dotenv"},
				},
			},
			{
				Name: "empty",
				Vars: map[string]string{"FOO": "bar"},
			},
		},
	}
}

func TestSync_DryRun(t *testing.T) {
	cfg := syncTestConfig()
	tmpDir := t.TempDir()

	results, err := Sync(cfg, "app", tmpDir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// Files must NOT be created in dry-run mode.
	for _, r := range results {
		if _, err := os.Stat(r.File); !os.IsNotExist(err) {
			t.Errorf("file %q should not exist in dry-run mode", r.File)
		}
	}
}

func TestSync_WritesFiles(t *testing.T) {
	cfg := syncTestConfig()
	tmpDir := t.TempDir()

	results, err := Sync(cfg, "app", tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if _, err := os.Stat(filepath.Join(tmpDir, filepath.Base(r.File))); err != nil {
			// file path already includes tmpDir via outDir override
			if _, err2 := os.Stat(r.File); err2 != nil {
				t.Errorf("expected file %q to exist: %v", r.File, err2)
			}
		}
	}
}

func TestSync_NoTargets(t *testing.T) {
	cfg := syncTestConfig()
	_, err := Sync(cfg, "empty", "", false)
	if err == nil {
		t.Fatal("expected error for env set with no targets")
	}
}

func TestSync_UnknownSet(t *testing.T) {
	cfg := syncTestConfig()
	_, err := Sync(cfg, "nope", "", false)
	if err == nil {
		t.Fatal("expected error for unknown env set")
	}
}
