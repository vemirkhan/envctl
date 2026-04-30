package env

import (
	"testing"

	"github.com/user/envctl/internal/config"
)

func mergeTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "base",
				Base: map[string]string{
					"APP_NAME": "myapp",
					"LOG_LEVEL": "info",
				},
			},
			{
				Name: "overrides",
				Base: map[string]string{
					"LOG_LEVEL": "debug",
					"DEBUG": "true",
				},
			},
			{
				Name: "shared",
				Base: map[string]string{
					"REGION": "us-east-1",
					"APP_NAME": "myapp",
				},
			},
		},
	}
}

func TestMerge_SingleSet(t *testing.T) {
	cfg := mergeTestConfig()
	res, err := Merge(cfg, []string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", res.Vars["APP_NAME"])
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_LaterSetWins(t *testing.T) {
	cfg := mergeTestConfig()
	res, err := Merge(cfg, []string{"base", "overrides"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", res.Vars["LOG_LEVEL"])
	}
	if res.Vars["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", res.Vars["APP_NAME"])
	}
	if res.Vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", res.Vars["DEBUG"])
	}
}

func TestMerge_ConflictsDetected(t *testing.T) {
	cfg := mergeTestConfig()
	res, err := Merge(cfg, []string{"base", "overrides"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "LOG_LEVEL" {
		t.Errorf("expected conflict on LOG_LEVEL, got %q", res.Conflicts[0].Key)
	}
}

func TestMerge_NoConflictSameValue(t *testing.T) {
	cfg := mergeTestConfig()
	res, err := Merge(cfg, []string{"base", "shared"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// APP_NAME is the same in both sets — should not be a conflict
	for _, c := range res.Conflicts {
		if c.Key == "APP_NAME" {
			t.Errorf("APP_NAME should not be a conflict when values are identical")
		}
	}
}

func TestMerge_UnknownSet(t *testing.T) {
	cfg := mergeTestConfig()
	_, err := Merge(cfg, []string{"base", "nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestMerge_EmptySets(t *testing.T) {
	cfg := mergeTestConfig()
	_, err := Merge(cfg, []string{})
	if err == nil {
		t.Fatal("expected error for empty set list")
	}
}
