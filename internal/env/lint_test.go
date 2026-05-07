package env_test

import (
	"bytes"
	"testing"

	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

func lintTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "clean",
				Base: map[string]string{
					"APP_ENV": "production",
					"LOG_LEVEL": "info",
				},
				Targets: []config.Target{
					{Name: "staging", Overrides: map[string]string{"APP_ENV": "staging"}},
				},
			},
			{
				Name: "dirty",
				Base: map[string]string{
					"app_env":  "production",
					"LOG_LEVEL": "",
				},
				Targets: []config.Target{
					{Name: "staging", Overrides: map[string]string{"UNKNOWN_KEY": "val"}},
				},
			},
		},
	}
}

func TestLint_CleanSet(t *testing.T) {
	cfg := lintTestConfig()
	result, err := env.Lint(cfg, "clean")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK() {
		t.Errorf("expected no issues, got %d: %+v", len(result.Issues), result.Issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	cfg := lintTestConfig()
	result, err := env.Lint(cfg, "dirty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "app_env" {
			found = true
		}
	}
	if !found {
		t.Error("expected issue for lowercase key 'app_env'")
	}
}

func TestLint_EmptyValue(t *testing.T) {
	cfg := lintTestConfig()
	result, err := env.Lint(cfg, "dirty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "LOG_LEVEL" && issue.Message == "value is empty or whitespace-only" {
			found = true
		}
	}
	if !found {
		t.Error("expected issue for empty value on LOG_LEVEL")
	}
}

func TestLint_OverrideKeyNotInBase(t *testing.T) {
	cfg := lintTestConfig()
	result, err := env.Lint(cfg, "dirty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "UNKNOWN_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected issue for override key not in base")
	}
}

func TestLint_UnknownSet(t *testing.T) {
	cfg := lintTestConfig()
	_, err := env.Lint(cfg, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown set, got nil")
	}
}

func TestWriteLint_OK(t *testing.T) {
	result := &env.LintResult{}
	var buf bytes.Buffer
	env.WriteLint(&buf, result, "clean")
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("no issues")) {
		t.Errorf("expected 'no issues' in output, got: %s", buf.String())
	}
}
