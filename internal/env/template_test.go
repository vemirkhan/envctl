package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func templateTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"HOST": "localhost",
					"PORT": "8080",
					"DB":   "mydb",
				},
				Targets: map[string]map[string]string{
					"prod": {
						"HOST": "prod.example.com",
						"PORT": "443",
					},
				},
			},
		},
	}
}

func TestTemplate_BasicSubstitution(t *testing.T) {
	cfg := templateTestConfig()
	res, err := Template(cfg, "app", "", "Connect to {{HOST}}:{{PORT}}/{{DB}}", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "Connect to localhost:8080/mydb" {
		t.Errorf("got %q", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestTemplate_WithTarget(t *testing.T) {
	cfg := templateTestConfig()
	res, err := Template(cfg, "app", "prod", "{{HOST}}:{{PORT}}", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "prod.example.com:443" {
		t.Errorf("got %q", res.Rendered)
	}
}

func TestTemplate_MissingKeyLenient(t *testing.T) {
	cfg := templateTestConfig()
	res, err := Template(cfg, "app", "", "{{HOST}} and {{UNKNOWN}}", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "UNKNOWN" {
		t.Errorf("expected [UNKNOWN], got %v", res.Missing)
	}
	if res.Rendered != "localhost and {{UNKNOWN}}" {
		t.Errorf("got %q", res.Rendered)
	}
}

func TestTemplate_MissingKeyStrict(t *testing.T) {
	cfg := templateTestConfig()
	_, err := Template(cfg, "app", "", "{{HOST}} and {{MISSING}}", true)
	if err == nil {
		t.Fatal("expected error for strict mode with missing key")
	}
}

func TestTemplate_UnknownSet(t *testing.T) {
	cfg := templateTestConfig()
	_, err := Template(cfg, "ghost", "", "{{HOST}}", false)
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestTemplate_DuplicatePlaceholder(t *testing.T) {
	cfg := templateTestConfig()
	res, err := Template(cfg, "app", "", "{{NOPE}} and {{NOPE}} again", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing key, got %v", res.Missing)
	}
}
