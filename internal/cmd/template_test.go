package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envctl/envctl/internal/config"
	"gopkg.in/yaml.v3"
)

func writeTemplateConfig(t *testing.T) string {
	t.Helper()
	cfg := config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"HOST": "localhost",
					"PORT": "8080",
				},
				Targets: map[string]map[string]string{
					"prod": {"HOST": "prod.example.com"},
				},
			},
		},
	}
	data, _ := yaml.Marshal(&cfg)
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	os.WriteFile(p, data, 0644)
	return p
}

func runTemplateCmd(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	root := NewRootCmd()
	root.AddCommand(NewTemplateCmd())
	var out, errOut bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), errOut.String(), err
}

func TestTemplateCmd_InlineTemplate(t *testing.T) {
	cp := writeTemplateConfig(t)
	out, _, err := runTemplateCmd(t, "--config", cp, "template", "app", "{{HOST}}:{{PORT}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "localhost:8080") {
		t.Errorf("got %q", out)
	}
}

func TestTemplateCmd_WithTarget(t *testing.T) {
	cp := writeTemplateConfig(t)
	out, _, err := runTemplateCmd(t, "--config", cp, "template", "app", "{{HOST}}:{{PORT}}", "--target", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod.example.com:8080") {
		t.Errorf("got %q", out)
	}
}

func TestTemplateCmd_MissingKeyWarning(t *testing.T) {
	cp := writeTemplateConfig(t)
	_, errOut, err := runTemplateCmd(t, "--config", cp, "template", "app", "{{HOST}} {{MISSING}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(errOut, "MISSING") {
		t.Errorf("expected warning about MISSING, got stderr: %q", errOut)
	}
}

func TestTemplateCmd_StrictFails(t *testing.T) {
	cp := writeTemplateConfig(t)
	_, _, err := runTemplateCmd(t, "--config", cp, "template", "app", "{{NOPE}}", "--strict")
	if err == nil {
		t.Fatal("expected error in strict mode with unresolved placeholder")
	}
}

func TestTemplateCmd_MissingArgs(t *testing.T) {
	cp := writeTemplateConfig(t)
	_, _, err := runTemplateCmd(t, "--config", cp, "template")
	if err == nil {
		t.Fatal("expected error with no args")
	}
}
