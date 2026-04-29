package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func validationTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Vars: map[string]string{
					"APP_HOST": "localhost",
					"APP_PORT": "8080",
				},
				Targets: []config.Target{
					{
						Name:      "prod",
						Overrides: map[string]string{"APP_HOST": "prod.example.com"},
					},
				},
			},
		},
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := validationTestConfig()
	if err := Validate(cfg, "app"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_UnknownSet(t *testing.T) {
	cfg := validationTestConfig()
	if err := Validate(cfg, "missing"); err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestValidate_InvalidBaseKey(t *testing.T) {
	cfg := validationTestConfig()
	cfg.EnvSets[0].Vars["123INVALID"] = "value"
	err := Validate(cfg, "app")
	if err == nil {
		t.Fatal("expected validation error for invalid key")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) == 0 {
		t.Fatal("expected at least one issue")
	}
}

func TestValidate_EmptyBaseValue(t *testing.T) {
	cfg := validationTestConfig()
	cfg.EnvSets[0].Vars["APP_HOST"] = "   "
	err := Validate(cfg, "app")
	if err == nil {
		t.Fatal("expected validation error for empty value")
	}
}

func TestValidate_EmptyTargetOverride(t *testing.T) {
	cfg := validationTestConfig()
	cfg.EnvSets[0].Targets[0].Overrides["APP_HOST"] = ""
	err := Validate(cfg, "app")
	if err == nil {
		t.Fatal("expected validation error for empty target override")
	}
}

func TestValidate_InvalidTargetKey(t *testing.T) {
	cfg := validationTestConfig()
	cfg.EnvSets[0].Targets[0].Overrides["bad-key"] = "value"
	err := Validate(cfg, "app")
	if err == nil {
		t.Fatal("expected validation error for invalid target key")
	}
}
