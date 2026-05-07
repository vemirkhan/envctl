package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func writeMaskConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "envctl.yaml")
	data := `env_sets:
  - name: production
    base:
      DB_PASSWORD: secret
      API_KEY: key123
      APP_ENV: prod
    sealed:
      - DB_PASSWORD
`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}

func runMaskCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd()
	root.AddCommand(NewMaskCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestMaskCmd_AddKey(t *testing.T) {
	path := writeMaskConfig(t)
	out, err := runMaskCmd(t, "--config", path, "mask", "add", "production", "API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "masked: API_KEY") {
		t.Errorf("expected masked output, got: %q", out)
	}

	cfg, _ := config.Load(path)
	set := cfg.EnvSetByName("production")
	found := false
	for _, k := range set.Sealed {
		if k == "API_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected API_KEY to be sealed after mask add")
	}
}

func TestMaskCmd_AddAlreadySealed(t *testing.T) {
	path := writeMaskConfig(t)
	out, err := runMaskCmd(t, "--config", path, "mask", "add", "production", "DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "already masked") {
		t.Errorf("expected already masked message, got: %q", out)
	}
}

func TestMaskCmd_RemoveKey(t *testing.T) {
	path := writeMaskConfig(t)
	out, err := runMaskCmd(t, "--config", path, "mask", "remove", "production", "DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "unmasked: DB_PASSWORD") {
		t.Errorf("expected unmasked output, got: %q", out)
	}

	cfg, _ := config.Load(path)
	set := cfg.EnvSetByName("production")
	for _, k := range set.Sealed {
		if k == "DB_PASSWORD" {
			t.Error("expected DB_PASSWORD to be removed from sealed")
		}
	}
}

func TestMaskCmd_UnknownSet(t *testing.T) {
	path := writeMaskConfig(t)
	_, err := runMaskCmd(t, "--config", path, "mask", "add", "staging", "API_KEY")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestMaskCmd_UnknownKey(t *testing.T) {
	path := writeMaskConfig(t)
	_, err := runMaskCmd(t, "--config", path, "mask", "add", "production", "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}
