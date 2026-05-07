package env

import (
	"os"
	"path/filepath"
	"testing"

	"envctl/internal/config"
)

func importTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"EXISTING": "old",
				},
				Targets: map[string]map[string]string{
					"prod": {"PROD_ONLY": "yes"},
				},
			},
		},
	}
}

func writeTempImportFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "import.env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestImport_DotenvBase(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, "NEW_KEY=hello\nANOTHER=world\n")

	n, err := Import(cfg, ImportOptions{File: path, Format: ImportFormatDotenv, SetName: "app", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 imported, got %d", n)
	}
	if cfg.EnvSets[0].Base["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello")
	}
}

func TestImport_NoOverwrite(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, "EXISTING=new_value\n")

	n, err := Import(cfg, ImportOptions{File: path, Format: ImportFormatDotenv, SetName: "app", Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 imported (no overwrite), got %d", n)
	}
	if cfg.EnvSets[0].Base["EXISTING"] != "old" {
		t.Errorf("expected EXISTING to remain 'old'")
	}
}

func TestImport_JSONFormat(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, `{"JSON_KEY":"json_val"}`)

	n, err := Import(cfg, ImportOptions{File: path, Format: ImportFormatJSON, SetName: "app", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 imported, got %d", n)
	}
	if cfg.EnvSets[0].Base["JSON_KEY"] != "json_val" {
		t.Errorf("expected JSON_KEY=json_val")
	}
}

func TestImport_IntoTarget(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, "TARGET_KEY=tval\n")

	_, err := Import(cfg, ImportOptions{File: path, Format: ImportFormatDotenv, SetName: "app", Target: "prod", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvSets[0].Targets["prod"]["TARGET_KEY"] != "tval" {
		t.Errorf("expected TARGET_KEY in prod target")
	}
}

func TestImport_UnknownSet(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, "KEY=val\n")

	_, err := Import(cfg, ImportOptions{File: path, Format: ImportFormatDotenv, SetName: "missing", Overwrite: true})
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestImport_UnsupportedFormat(t *testing.T) {
	cfg := importTestConfig()
	path := writeTempImportFile(t, "KEY=val\n")

	_, err := Import(cfg, ImportOptions{File: path, Format: "yaml", SetName: "app", Overwrite: true})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
