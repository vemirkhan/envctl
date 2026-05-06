package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePromoteConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

const promoteYAML = `
env_sets:
  - name: app
    base:
      HOST: localhost
      PORT: "8080"
    targets:
      staging:
        HOST: staging.example.com
        DEBUG: "true"
      production:
        HOST: prod.example.com
`

func TestPromoteCmd_Success(t *testing.T) {
	cfgPath := writePromoteConfig(t, promoteYAML)
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--config", cfgPath, "promote", "app", "staging", "production"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DEBUG") {
		t.Errorf("expected DEBUG in output, got: %s", out)
	}
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected target names in output, got: %s", out)
	}
}

func TestPromoteCmd_NothingToPromote(t *testing.T) {
	cfgPath := writePromoteConfig(t, promoteYAML)
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	// promote production -> staging: HOST already exists in staging, no new keys
	root.SetArgs([]string{"--config", cfgPath, "promote", "app", "production", "staging"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No keys promoted") {
		t.Errorf("expected 'No keys promoted' message, got: %s", buf.String())
	}
}

func TestPromoteCmd_UnknownSet(t *testing.T) {
	cfgPath := writePromoteConfig(t, promoteYAML)
	root := NewRootCmd()
	root.SetArgs([]string{"--config", cfgPath, "promote", "missing", "staging", "production"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for unknown env set")
	}
}

func TestPromoteCmd_MissingArgs(t *testing.T) {
	cfgPath := writePromoteConfig(t, promoteYAML)
	root := NewRootCmd()
	root.SetArgs([]string{"--config", cfgPath, "promote", "app", "staging"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing argument")
	}
}
