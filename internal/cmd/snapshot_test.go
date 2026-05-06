package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envctl/envctl/internal/cmd"
)

const snapshotConfigYAML = `
env_sets:
  - name: web
    base:
      HOST: localhost
      PORT: "3000"
    targets:
      prod:
        PORT: "443"
`

func writeSnapshotConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func runSnapshotCmd(t *testing.T, cfgPath string, args ...string) (string, error) {
	t.Helper()
	root := cmd.NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	fullArgs := append([]string{"--config", cfgPath, "snapshot"}, args...)
	root.SetArgs(fullArgs)
	err := root.Execute()
	return buf.String(), err
}

func TestSnapshotCmd_Take(t *testing.T) {
	p := writeSnapshotConfig(t, snapshotConfigYAML)
	out, err := runSnapshotCmd(t, p, "take", "web", "--name", "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "snap1") {
		t.Errorf("expected snap1 in output, got: %s", out)
	}
}

func TestSnapshotCmd_TakeWithTarget(t *testing.T) {
	p := writeSnapshotConfig(t, snapshotConfigYAML)
	out, err := runSnapshotCmd(t, p, "take", "web", "--target", "prod", "--name", "prod-snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod-snap") {
		t.Errorf("expected prod-snap in output, got: %s", out)
	}
}

func TestSnapshotCmd_TakeUnknownSet(t *testing.T) {
	p := writeSnapshotConfig(t, snapshotConfigYAML)
	_, err := runSnapshotCmd(t, p, "take", "ghost", "--name", "x")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestSnapshotCmd_List(t *testing.T) {
	p := writeSnapshotConfig(t, snapshotConfigYAML)
	// take a snapshot first so list has something
	_, _ = runSnapshotCmd(t, p, "take", "web", "--name", "list-test")
	out, err := runSnapshotCmd(t, p, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "list-test") {
		t.Errorf("expected list-test in list output, got: %s", out)
	}
}

func TestSnapshotCmd_Delete(t *testing.T) {
	p := writeSnapshotConfig(t, snapshotConfigYAML)
	_, _ = runSnapshotCmd(t, p, "take", "web", "--name", "del-me")
	out, err := runSnapshotCmd(t, p, "delete", "del-me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected deleted in output, got: %s", out)
	}
}
