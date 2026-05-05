package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envctl/internal/cmd"
)

func writeListConfig(t *testing.T) string {
	t.Helper()
	content := `
env_sets:
  - name: alpha
    base:
      APP: alpha
      PORT: "3000"
    targets:
      - name: prod
        overrides:
          PORT: "4000"
  - name: beta
    base:
      APP: beta
`
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return p
}

func runListCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmdForTest()
	root.AddCommand(cmd.NewListCmd())
	return executeCmd(root, args...)
}

func NewRootCmdForTest() *cobra.Command {
	return cmd.NewRootCmd()
}

func TestListCmd_ShowsAllSets(t *testing.T) {
	cfgPath := writeListConfig(t)
	out, err := runListCmd(t, "--config", cfgPath, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output, got: %s", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected 'beta' in output, got: %s", out)
	}
}

func TestListCmd_ShowsTargets(t *testing.T) {
	cfgPath := writeListConfig(t)
	out, err := runListCmd(t, "--config", cfgPath, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected 'prod' target in output, got: %s", out)
	}
}

func TestListCmd_InvalidConfig(t *testing.T) {
	_, err := runListCmd(t, "--config", "/nonexistent/path.yaml", "list")
	if err == nil {
		t.Error("expected error for missing config, got nil")
	}
}

func TestListCmd_ShowsBaseVarCount(t *testing.T) {
	cfgPath := writeListConfig(t)
	out, err := runListCmd(t, "--config", cfgPath, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// alpha has 2 base vars
	if !strings.Contains(out, "2") {
		t.Errorf("expected base var count in output, got: %s", out)
	}
}
