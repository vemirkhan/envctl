package env

import (
	"strings"
	"testing"
)

var sampleVars = map[string]string{
	"APP_ENV":   "production",
	"DB_URL":    "postgres://localhost/mydb",
	"SECRET_KEY": "s3cr3t value",
}

func TestExport_FormatExport(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleVars, FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export SECRET_KEY=") {
		t.Errorf("expected export SECRET_KEY= in output, got:\n%s", out)
	}
}

func TestExport_FormatDotenv(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleVars, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("dotenv format should not contain 'export', got:\n%s", out)
	}
}

func TestExport_FormatJSON(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleVars, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("expected APP_ENV key in JSON output, got:\n%s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleVars, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected first line to be APP_ENV, got: %s", lines[0])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var sb strings.Builder
	err := Export(&sb, sampleVars, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_EmptyVars(t *testing.T) {
	for _, format := range []Format{FormatExport, FormatDotenv, FormatJSON} {
		var sb strings.Builder
		if err := Export(&sb, map[string]string{}, format); err != nil {
			t.Errorf("format %q: unexpected error for empty vars: %v", format, err)
		}
	}
}
