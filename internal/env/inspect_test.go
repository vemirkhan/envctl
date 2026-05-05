package env

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envctl/internal/config"
)

func inspectTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"APP_ENV": "production",
					"LOG_LEVEL": "warn",
				},
				Targets: []config.Target{
					{
						Name: "us-east",
						Overrides: map[string]string{"REGION": "us-east-1"},
					},
					{
						Name: "eu-west",
						Overrides: map[string]string{"REGION": "eu-west-1"},
					},
				},
			},
			{
				Name: "staging",
				Base: map[string]string{"APP_ENV": "staging"},
			},
		},
	}
}

func TestInspect_KnownSet(t *testing.T) {
	cfg := inspectTestConfig()
	r, err := Inspect(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "production" {
		t.Errorf("expected name production, got %s", r.Name)
	}
	if r.Base["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", r.Base["APP_ENV"])
	}
	if len(r.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(r.Targets))
	}
}

func TestInspect_NoTargets(t *testing.T) {
	cfg := inspectTestConfig()
	r, err := Inspect(cfg, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Targets) != 0 {
		t.Errorf("expected 0 targets, got %d", len(r.Targets))
	}
}

func TestInspect_UnknownSet(t *testing.T) {
	cfg := inspectTestConfig()
	_, err := Inspect(cfg, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestWriteInspect_Output(t *testing.T) {
	cfg := inspectTestConfig()
	r, _ := Inspect(cfg, "production")
	var buf bytes.Buffer
	WriteInspect(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "Env Set: production") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Error("expected APP_ENV in base output")
	}
	if !strings.Contains(out, "Target: eu-west") {
		t.Error("expected eu-west target in output")
	}
	if !strings.Contains(out, "REGION=us-east-1") {
		t.Error("expected us-east-1 override in output")
	}
}
