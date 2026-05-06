package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTagConfig(t *testing.T) string {
	t.Helper()
	content := `env_sets:
  - name: web
    tags: [prod]
    base:
      PORT: "8080"
  - name: worker
    tags: []
    base:
      QUEUE: "default"
`
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func runTagCmd(t *testing.T, cfgPath string, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd()
	out := &strings.Builder{}
	root.SetOut(out)
	root.SetErr(out)
	fullArgs := append([]string{"--config", cfgPath, "tag"}, args...)
	root.SetArgs(fullArgs)
	err := root.Execute()
	return out.String(), err
}

func TestTagCmd_AddTag(t *testing.T) {
	p := writeTagConfig(t)
	out, err := runTagCmd(t, p, "add", "worker", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "tagged") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestTagCmd_RemoveTag(t *testing.T) {
	p := writeTagConfig(t)
	out, err := runTagCmd(t, p, "remove", "web", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestTagCmd_ListByTag(t *testing.T) {
	p := writeTagConfig(t)
	out, err := runTagCmd(t, p, "list", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected 'web' in output, got: %s", out)
	}
}

func TestTagCmd_ListByTag_NoMatches(t *testing.T) {
	p := writeTagConfig(t)
	out, err := runTagCmd(t, p, "list", "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no env sets") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestTagCmd_AddUnknownSet(t *testing.T) {
	p := writeTagConfig(t)
	_, err := runTagCmd(t, p, "add", "ghost", "x")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}
