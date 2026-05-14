package env

import (
	"bytes"
	"testing"

	"github.com/user/envctl/internal/config"
)

func grepTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []*config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"DATABASE_URL": "postgres://localhost/dev",
					"API_KEY":      "abc123",
					"LOG_LEVEL":    "debug",
				},
				Targets: []config.Target{
					{
						Name: "prod",
						Overrides: map[string]string{
							"DATABASE_URL": "postgres://prod-host/app",
							"LOG_LEVEL":    "warn",
						},
					},
				},
			},
			{
				Name: "worker",
				Base: map[string]string{
					"QUEUE_URL": "redis://localhost",
					"API_KEY":   "xyz789",
				},
			},
		},
	}
}

func TestGrep_MatchKey(t *testing.T) {
	cfg := grepTestConfig()
	results, err := Grep(cfg, "", "API_KEY", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestGrep_MatchValue(t *testing.T) {
	cfg := grepTestConfig()
	results, err := Grep(cfg, "", "localhost", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	for _, r := range results {
		if r.Value == "" {
			t.Error("result has empty value")
		}
	}
}

func TestGrep_ScopedToSet(t *testing.T) {
	cfg := grepTestConfig()
	results, err := Grep(cfg, "worker", "API_KEY", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].SetName != "worker" {
		t.Errorf("expected set 'worker', got %q", results[0].SetName)
	}
}

func TestGrep_UnknownSet(t *testing.T) {
	cfg := grepTestConfig()
	_, err := Grep(cfg, "missing", "KEY", false)
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestGrep_InvalidPattern(t *testing.T) {
	cfg := grepTestConfig()
	_, err := Grep(cfg, "", "[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestGrep_NoMatches(t *testing.T) {
	cfg := grepTestConfig()
	results, err := Grep(cfg, "", "NONEXISTENT_KEY_XYZ", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestWriteGrep_Output(t *testing.T) {
	results := []GrepResult{
		{SetName: "app", Target: "", Key: "API_KEY", Value: "abc123"},
		{SetName: "app", Target: "prod", Key: "LOG_LEVEL", Value: "warn"},
	}
	var buf bytes.Buffer
	WriteGrep(&buf, results)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("[app]")) {
		t.Error("expected output to contain '[app]'")
	}
	if !bytes.Contains(buf.Bytes(), []byte("[app:prod]")) {
		t.Error("expected output to contain '[app:prod]'")
	}
}

func TestWriteGrep_NoMatches(t *testing.T) {
	var buf bytes.Buffer
	WriteGrep(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("no matches")) {
		t.Error("expected 'no matches' message")
	}
}
