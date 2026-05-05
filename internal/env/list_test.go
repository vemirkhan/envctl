package env_test

import (
	"bytes"
	"strings"
	"testing"

	"envctl/internal/config"
	"envctl/internal/env"
)

func listTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{"APP_ENV": "prod", "PORT": "8080", "LOG": "info"},
				Targets: []config.Target{
					{Name: "us-east", Overrides: map[string]string{"PORT": "9090"}},
					{Name: "eu-west", Overrides: map[string]string{}},
				},
			},
			{
				Name: "staging",
				Base: map[string]string{"APP_ENV": "staging"},
				Targets: []config.Target{},
			},
		},
	}
}

func TestList_ReturnsAllSets(t *testing.T) {
	cfg := listTestConfig()
	results := env.List(cfg)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestList_SortedByName(t *testing.T) {
	cfg := listTestConfig()
	results := env.List(cfg)

	if results[0].Name != "production" {
		t.Errorf("expected first to be 'production', got %q", results[0].Name)
	}
	if results[1].Name != "staging" {
		t.Errorf("expected second to be 'staging', got %q", results[1].Name)
	}
}

func TestList_BaseLen(t *testing.T) {
	cfg := listTestConfig()
	results := env.List(cfg)

	for _, r := range results {
		if r.Name == "production" && r.BaseLen != 3 {
			t.Errorf("expected production BaseLen=3, got %d", r.BaseLen)
		}
		if r.Name == "staging" && r.BaseLen != 1 {
			t.Errorf("expected staging BaseLen=1, got %d", r.BaseLen)
		}
	}
}

func TestList_TargetsSorted(t *testing.T) {
	cfg := listTestConfig()
	results := env.List(cfg)

	for _, r := range results {
		if r.Name == "production" {
			if len(r.Targets) != 2 {
				t.Fatalf("expected 2 targets, got %d", len(r.Targets))
			}
			if r.Targets[0] != "eu-west" || r.Targets[1] != "us-east" {
				t.Errorf("targets not sorted: %v", r.Targets)
			}
		}
	}
}

func TestWriteList_Empty(t *testing.T) {
	var buf bytes.Buffer
	env.WriteList(&buf, []env.ListResult{})
	if !strings.Contains(buf.String(), "No env sets defined") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestWriteList_Output(t *testing.T) {
	cfg := listTestConfig()
	results := env.List(cfg)
	var buf bytes.Buffer
	env.WriteList(&buf, results)
	out := buf.String()

	if !strings.Contains(out, "production") {
		t.Error("expected 'production' in output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected 'staging' in output")
	}
	if !strings.Contains(out, "eu-west") {
		t.Error("expected 'eu-west' in output")
	}
}
