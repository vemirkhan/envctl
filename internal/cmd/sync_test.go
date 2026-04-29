package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envctl/internal/config"
	"gopkg.in/yaml.v3"
)

func writeSyncConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "envctl.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshalling config: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	return path
}

func makeSyncCfg(outDir string) *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "web",
				Vars: map[string]string{"PORT": "8080", "HOST": "localhost"},
				Targets: []config.Target{
					{Name: "local", File: filepath.Join(outDir, "local.env"), Format: "dotenv"},
				},
			},
		},
	}
}

func TestSyncCmd_DryRun(t *testing.T) {
	outDir := t.TempDir()
	cfgPath := writeSyncConfig(t, makeSyncCfg(outDir))

	cmd := NewSyncCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--config", cfgPath, "--dry-run", "web"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Error("expected dry-run output, got none")
	}
	// No file should have been created.
	if _, err := os.Stat(filepath.Join(outDir, "local.env")); !os.IsNotExist(err) {
		t.Error("file should not exist after dry-run")
	}
}

func TestSyncCmd_WritesFile(t *testing.T) {
	outDir := t.TempDir()
	cfgPath := writeSyncConfig(t, makeSyncCfg(outDir))

	cmd := NewSyncCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--config", cfgPath, "web"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "local.env")); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

func TestSyncCmd_UnknownSet(t *testing.T) {
	outDir := t.TempDir()
	cfgPath := writeSyncConfig(t, makeSyncCfg(outDir))

	cmd := NewSyncCmd()
	cmd.SetArgs([]string{"--config", cfgPath, "nonexistent"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for unknown env set")
	}
}
