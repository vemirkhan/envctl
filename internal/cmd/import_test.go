package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"envctl/internal/config"
)

func writeImportConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envctl.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeImportSourceFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vars.env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runImportCmd(args ...string) (string, error) {
	root := NewRootCmd()
	root.AddCommand(NewImportCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestImportCmd_DotenvBase(t *testing.T) {
	cfg := &config.Config{
		EnvSets: []config.EnvSet{
			{Name: "app", Base: map[string]string{"OLD": "val"}, Targets: map[string]map[string]string{}},
		},
	}
	cfgPath := writeImportConfig(t, cfg)
	srcPath := writeImportSourceFile(t, "NEW_KEY=hello\n")

	out, err := runImportCmd("--config", cfgPath, "import", "app", srcPath, "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Imported 1") {
		t.Errorf("expected import count in output, got: %s", out)
	}
}

func TestImportCmd_WithTarget(t *testing.T) {
	cfg := &config.Config{
		EnvSets: []config.EnvSet{
			{Name: "app", Base: map[string]string{"A": "1"}, Targets: map[string]map[string]string{
				"staging": {},
			}},
		},
	}
	cfgPath := writeImportConfig(t, cfg)
	srcPath := writeImportSourceFile(t, "STAGE_VAR=stg\n")

	out, err := runImportCmd("--config", cfgPath, "import", "app", srcPath, "--target", "staging", "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "target: staging") {
		t.Errorf("expected target label in output, got: %s", out)
	}
}

func TestImportCmd_UnknownSet(t *testing.T) {
	cfg := &config.Config{
		EnvSets: []config.EnvSet{
			{Name: "app", Base: map[string]string{}, Targets: map[string]map[string]string{}},
		},
	}
	cfgPath := writeImportConfig(t, cfg)
	srcPath := writeImportSourceFile(t, "KEY=val\n")

	_, err := runImportCmd("--config", cfgPath, "import", "missing", srcPath)
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestImportCmd_MissingArgs(t *testing.T) {
	_, err := runImportCmd("import", "only-one-arg")
	if err == nil {
		t.Error("expected error for missing file argument")
	}
}
