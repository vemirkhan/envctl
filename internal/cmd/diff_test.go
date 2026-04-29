package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeDiffConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeDiffConfig: %v", err)
	}
	return p
}

func runDiffCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envctl"}
	root.AddCommand(NewDiffCmd())
	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

const diffCfg = `
env_sets:
  - name: production
    base:
      APP_ENV: production
      DB_HOST: db.prod.example.com
      LOG_LEVEL: info
    targets:
      - name: us-east
        overrides:
          DB_HOST: db.us-east.example.com
          LOG_LEVEL: warn
`

func TestDiffCmd_BaseVsTarget(t *testing.T) {
	p := writeDiffConfig(t, diffCfg)
	out, err := runDiffCmd(t, []string{"diff", "--config", p, "production", "us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in diff output, got: %s", out)
	}
}

func TestDiffCmd_UnknownSet(t *testing.T) {
	p := writeDiffConfig(t, diffCfg)
	_, err := runDiffCmd(t, []string{"diff", "--config", p, "staging", "us-east"})
	if err == nil {
		t.Fatal("expected error for unknown env set")
	}
}

func TestDiffCmd_UnknownTarget(t *testing.T) {
	p := writeDiffConfig(t, diffCfg)
	_, err := runDiffCmd(t, []string{"diff", "--config", p, "production", "eu-west"})
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestDiffCmd_MissingArgs(t *testing.T) {
	p := writeDiffConfig(t, diffCfg)
	_, err := runDiffCmd(t, []string{"diff", "--config", p, "production"})
	if err == nil {
		t.Fatal("expected error when target arg is missing")
	}
}
