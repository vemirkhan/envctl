package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSetkeyConfig(t *testing.T, dir string) string {
	t.Helper()
	content := `envsets:
  - name: app
    base:
      EXISTING: "yes"
    targets:
      - name: prod
        overrides:
          PROD_ONLY: "1"
`
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runSetkeyCmd(args ...string) (string, error) {
	root := NewRootCmd()
	root.AddCommand(NewSetKeyCmd())
	buf := new(strings.Builder)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestSetkeyCmd_NewBaseKey(t *testing.T) {
	dir := t.TempDir()
	cp := writeSetkeyConfig(t, dir)
	out, err := runSetkeyCmd("--config", cp, "setkey", "app", "NEW_KEY", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "NEW_KEY=hello") {
		t.Errorf("expected output to mention NEW_KEY=hello, got: %s", out)
	}
}

func TestSetkeyCmd_OverwriteFlag(t *testing.T) {
	dir := t.TempDir()
	cp := writeSetkeyConfig(t, dir)
	_, err := runSetkeyCmd("--config", cp, "setkey", "--overwrite", "app", "EXISTING", "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetkeyCmd_ConflictNoOverwrite(t *testing.T) {
	dir := t.TempDir()
	cp := writeSetkeyConfig(t, dir)
	_, err := runSetkeyCmd("--config", cp, "setkey", "app", "EXISTING", "bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetkeyCmd_TargetFlag(t *testing.T) {
	dir := t.TempDir()
	cp := writeSetkeyConfig(t, dir)
	out, err := runSetkeyCmd("--config", cp, "setkey", "--target", "prod", "app", "NEW_PROD", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected output to mention prod target, got: %s", out)
	}
}

func TestSetkeyCmd_UnknownSet(t *testing.T) {
	dir := t.TempDir()
	cp := writeSetkeyConfig(t, dir)
	_, err := runSetkeyCmd("--config", cp, "setkey", "ghost", "K", "v")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}
