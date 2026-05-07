package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envctl/internal/config"
)

func writeRollbackConfig(t *testing.T) string {
	t.Helper()
	content := `
env_sets:
  - name: production
    base:
      APP_ENV: prod
      LOG_LEVEL: warn
    targets:
      - name: us-east
        overrides:
          REGION: us-east-1
snapshots:
  - name: v1
    env_set: production
    vars:
      APP_ENV: staging
      LOG_LEVEL: debug
    targets:
      us-east:
        REGION: us-east-2
`
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runRollbackCmd(t *testing.T, cfgPath string, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd()
	root.AddCommand(NewRollbackCmd())
	var buf bytes.Buffer
	root.SetOut(&buf)
	cmdArgs := append([]string{"--config", cfgPath, "rollback"}, args...)
	root.SetArgs(cmdArgs)
	err := root.Execute()
	return buf.String(), err
}

func TestRollbackCmd_DryRun(t *testing.T) {
	cfgPath := writeRollbackConfig(t)
	out, err := runRollbackCmd(t, cfgPath, "--dry-run", "production", "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run notice, got: %s", out)
	}
	if !strings.Contains(out, "APP_ENV=staging") {
		t.Errorf("expected APP_ENV in output, got: %s", out)
	}
	// Config should be unchanged on disk (still has prod value)
	cfg, _ := config.Load(cfgPath)
	set := cfg.EnvSetByName("production")
	if set.Base["APP_ENV"] != "prod" {
		t.Errorf("dry-run should not have modified config")
	}
}

func TestRollbackCmd_RestoresBase(t *testing.T) {
	cfgPath := writeRollbackConfig(t)
	out, err := runRollbackCmd(t, cfgPath, "production", "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2 key(s) restored") {
		t.Errorf("expected restored count in output, got: %s", out)
	}
}

func TestRollbackCmd_UnknownSnapshot(t *testing.T) {
	cfgPath := writeRollbackConfig(t)
	_, err := runRollbackCmd(t, cfgPath, "production", "ghost")
	if err == nil {
		t.Error("expected error for unknown snapshot")
	}
}

func TestRollbackCmd_MissingArgs(t *testing.T) {
	cfgPath := writeRollbackConfig(t)
	_, err := runRollbackCmd(t, cfgPath, "production")
	if err == nil {
		t.Error("expected error for missing snapshot arg")
	}
}
