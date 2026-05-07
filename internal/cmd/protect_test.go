package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProtectConfig(t *testing.T) string {
	t.Helper()
	content := `env_sets:
  - name: production
    base:
      DB_URL: postgres://prod
      API_KEY: secret
      LOG_LEVEL: warn
    protected: []
`
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runProtectCmd(t *testing.T, cfgPath string, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd()
	all := append([]string{"--config", cfgPath, "protect"}, args...)
	out, err := executeCmd(root, all...)
	return out, err
}

func TestProtectCmd_AddKey(t *testing.T) {
	p := writeProtectConfig(t)
	_, err := runProtectCmd(t, p, "add", "production", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProtectCmd_ListKeys(t *testing.T) {
	p := writeProtectConfig(t)
	_, _ = runProtectCmd(t, p, "add", "production", "API_KEY")
	out, err := runProtectCmd(t, p, "list", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestProtectCmd_RemoveKey(t *testing.T) {
	p := writeProtectConfig(t)
	_, _ = runProtectCmd(t, p, "add", "production", "DB_URL", "API_KEY")
	_, err := runProtectCmd(t, p, "remove", "production", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := runProtectCmd(t, p, "list", "production")
	if strings.Contains(out, "DB_URL") {
		t.Errorf("DB_URL should have been removed, got: %s", out)
	}
}

func TestProtectCmd_UnknownSet(t *testing.T) {
	p := writeProtectConfig(t)
	_, err := runProtectCmd(t, p, "add", "staging", "DB_URL")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestProtectCmd_ListEmpty(t *testing.T) {
	p := writeProtectConfig(t)
	out, err := runProtectCmd(t, p, "list", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no protected keys") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
