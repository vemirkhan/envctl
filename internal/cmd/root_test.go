package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// executeCmd is a shared helper that runs a cobra command with the given args
// and returns combined stdout output and any error.
func executeCmd(cmd *cobra.Command, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestRootCmd_NoArgs(t *testing.T) {
	root := NewRootCmdForTest()
	out, err := executeCmd(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "envctl") {
		t.Errorf("expected help output to mention 'envctl', got: %q", out)
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	root := NewRootCmdForTest()
	names := make(map[string]bool)
	for _, sub := range root.Commands() {
		names[sub.Name()] = true
	}

	expected := []string{"export", "diff", "sync", "validate"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestRootCmd_ConfigFlag(t *testing.T) {
	root := NewRootCmdForTest()
	flag := root.PersistentFlags().Lookup("config")
	if flag == nil {
		t.Fatal("expected --config persistent flag to exist")
	}
	if flag.DefValue != "envctl.yaml" {
		t.Errorf("expected default config value 'envctl.yaml', got %q", flag.DefValue)
	}
}

func TestRootCmd_ShortConfigFlag(t *testing.T) {
	root := NewRootCmdForTest()
	flag := root.PersistentFlags().ShorthandLookup("c")
	if flag == nil {
		t.Fatal("expected -c shorthand flag to exist")
	}
}
