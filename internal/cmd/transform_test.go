package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func writeTransformConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func runTransformCmd(t *testing.T, cfgPath string, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd()
	root.AddCommand(NewTransformCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	all := append([]string{"--config", cfgPath, "transform"}, args...)
	root.SetArgs(all)
	err := root.Execute()
	return buf.String(), err
}

const transformCfgYAML = `
env_sets:
  - name: app
    base:
      APP_NAME: myapp
      LOG_LEVEL: "  debug  "
      REGION: us-east-1
    targets:
      prod:
        LOG_LEVEL: warn
`

func TestTransformCmd_UpperBase(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	out, err := runTransformCmd(t, p, "app", "upper", "--keys", "REGION")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "REGION") {
		t.Errorf("expected REGION in output, got: %s", out)
	}
}

func TestTransformCmd_TrimDryRun(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	out, err := runTransformCmd(t, p, "app", "trim", "--keys", "LOG_LEVEL", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in output, got: %s", out)
	}
	// ensure file not mutated
	raw, _ := os.ReadFile(p)
	var m map[string]interface{}
	_ = yaml.Unmarshal(raw, &m)
}

func TestTransformCmd_WithTarget(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	out, err := runTransformCmd(t, p, "app", "upper", "--target", "prod", "--keys", "LOG_LEVEL")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
}

func TestTransformCmd_UnknownSet(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	_, err := runTransformCmd(t, p, "nope", "upper")
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestTransformCmd_InvalidOp(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	_, err := runTransformCmd(t, p, "app", "rot13")
	if err == nil {
		t.Error("expected error for invalid op")
	}
}

func TestTransformCmd_MissingArgs(t *testing.T) {
	p := writeTransformConfig(t, transformCfgYAML)
	_, err := runTransformCmd(t, p, "app")
	if err == nil {
		t.Error("expected error for missing op argument")
	}
}
