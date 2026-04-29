package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeValidateConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeValidateConfig: %v", err)
	}
	return p
}

func runValidateCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envctl"}
	root.AddCommand(NewValidateCmd())
	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

const validateCfg = `
env_sets:
  - name: production
    base:
      APP_ENV: production
      DB_HOST: db.prod.example.com
    targets:
      - name: us-east
        overrides:
          DB_HOST: db.us-east.example.com
`

func TestValidateCmd_Valid(t *testing.T) {
	p := writeValidateConfig(t, validateCfg)
	out, err := runValidateCmd(t, []string{"validate", "--config", p, "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "valid") {
		t.Errorf("expected 'valid' in output, got: %s", out)
	}
}

func TestValidateCmd_UnknownSet(t *testing.T) {
	p := writeValidateConfig(t, validateCfg)
	_, err := runValidateCmd(t, []string{"validate", "--config", p, "staging"})
	if err == nil {
		t.Fatal("expected error for unknown env set")
	}
}

func TestValidateCmd_MissingArg(t *testing.T) {
	p := writeValidateConfig(t, validateCfg)
	_, err := runValidateCmd(t, []string{"validate", "--config", p})
	if err == nil {
		t.Fatal("expected error when env set name is missing")
	}
}

func TestValidateCmd_InvalidConfig(t *testing.T) {
	p := writeValidateConfig(t, `env_sets:\n  - name: bad\n    base:\n      ": broken`)
	_, err := runValidateCmd(t, []string{"validate", "--config", p, "bad"})
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
