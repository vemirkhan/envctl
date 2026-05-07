package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeReorderConfig(t *testing.T) string {
	t.Helper()
	content := `env_sets:
  - name: app
    base:
      HOST: localhost
      PORT: "8080"
      DEBUG: "false"
    key_order: [HOST, PORT, DEBUG]
`
	f, err := os.CreateTemp(t.TempDir(), "envctl-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestReorderCmd_Success(t *testing.T) {
	cfgPath := writeReorderConfig(t)
	root := NewRootCmd()
	out, err := executeCmd(root, "-c", cfgPath, "reorder", "app", "DEBUG,HOST,PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "reordered keys") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestReorderCmd_UnknownSet(t *testing.T) {
	cfgPath := writeReorderConfig(t)
	root := NewRootCmd()
	_, err := executeCmd(root, "-c", cfgPath, "reorder", "ghost", "HOST")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestReorderCmd_UnknownKey(t *testing.T) {
	cfgPath := writeReorderConfig(t)
	root := NewRootCmd()
	_, err := executeCmd(root, "-c", cfgPath, "reorder", "app", "MISSING")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestReorderCmd_MissingArgs(t *testing.T) {
	cfgPath := writeReorderConfig(t)
	root := NewRootCmd()
	_, err := executeCmd(root, "-c", cfgPath, "reorder", "app")
	if err == nil {
		t.Fatal("expected error for missing key argument")
	}
}
