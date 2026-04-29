package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envctl/internal/cmd"
)

func writeExportConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

const exportTestConfig = `
env_sets:
  - name: app
    vars:
      APP_ENV: development
      LOG_LEVEL: debug
    targets:
      production:
        vars:
          APP_ENV: production
          LOG_LEVEL: warn
`

func TestExportCmd_DefaultFormat(t *testing.T) {
	cfgPath := writeExportConfig(t, exportTestConfig)

	c := cmd.NewExportCmd()
	c.SetArgs([]string{"--config", cfgPath, "app"})

	buf := &bytes.Buffer{}
	c.SetOut(buf)

	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportCmd_DotenvFormat(t *testing.T) {
	cfgPath := writeExportConfig(t, exportTestConfig)

	c := cmd.NewExportCmd()
	c.SetArgs([]string{"--config", cfgPath, "--format", "dotenv", "app"})

	buf := &bytes.Buffer{}
	c.SetOut(buf)

	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportCmd_WithTarget(t *testing.T) {
	cfgPath := writeExportConfig(t, exportTestConfig)

	c := cmd.NewExportCmd()
	c.SetArgs([]string{"--config", cfgPath, "--target", "production", "--format", "json", "app"})

	buf := &bytes.Buffer{}
	c.SetOut(buf)

	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExportCmd_UnknownEnvSet(t *testing.T) {
	cfgPath := writeExportConfig(t, exportTestConfig)

	c := cmd.NewExportCmd()
	c.SetArgs([]string{"--config", cfgPath, "nonexistent"})

	err := c.Execute()
	if err == nil {
		t.Fatal("expected error for unknown env set, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("expected error to mention set name, got: %v", err)
	}
}

func TestExportCmd_MissingConfig(t *testing.T) {
	c := cmd.NewExportCmd()
	c.SetArgs([]string{"--config", "/nonexistent/path.yaml", "app"})

	if err := c.Execute(); err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}
