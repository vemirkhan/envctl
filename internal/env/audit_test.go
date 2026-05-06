package env

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func auditTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"HOST": "localhost",
					"PORT": "8080",
					"DEBUG": "",
				},
				Targets: []config.Target{
					{
						Name: "prod",
						Overrides: map[string]string{
							"HOST": "prod.example.com",
							"PORT": "8080", // same as base — redundant
							"EXTRA": "value", // not in base
						},
					},
				},
			},
		},
	}
}

func TestAudit_UnknownSet(t *testing.T) {
	cfg := auditTestConfig()
	_, err := Audit(cfg, "nope")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestAudit_EmptyBaseValue(t *testing.T) {
	cfg := auditTestConfig()
	results, err := Audit(cfg, "app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := results[0]
	if len(base.EmptyValues) != 1 || base.EmptyValues[0] != "DEBUG" {
		t.Errorf("expected DEBUG in empty values, got %v", base.EmptyValues)
	}
}

func TestAudit_RedundantOverride(t *testing.T) {
	cfg := auditTestConfig()
	results, err := Audit(cfg, "app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := results[1]
	if len(prod.UnusedOverrides) != 1 || prod.UnusedOverrides[0] != "PORT" {
		t.Errorf("expected PORT as redundant override, got %v", prod.UnusedOverrides)
	}
}

func TestAudit_MissingInBase(t *testing.T) {
	cfg := auditTestConfig()
	results, err := Audit(cfg, "app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prod := results[1]
	if len(prod.MissingInBase) != 1 || prod.MissingInBase[0] != "EXTRA" {
		t.Errorf("expected EXTRA as missing in base, got %v", prod.MissingInBase)
	}
}

func TestWriteAudit_Output(t *testing.T) {
	cfg := auditTestConfig()
	results, _ := Audit(cfg, "app")
	var buf bytes.Buffer
	WriteAudit(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "REDUNDANT override: PORT") {
		t.Errorf("expected redundant PORT in output, got:\n%s", out)
	}
	if !strings.Contains(out, "UNKNOWN key: EXTRA") {
		t.Errorf("expected unknown EXTRA in output, got:\n%s", out)
	}
	if !strings.Contains(out, "EMPTY value: DEBUG") {
		t.Errorf("expected empty DEBUG in output, got:\n%s", out)
	}
}
