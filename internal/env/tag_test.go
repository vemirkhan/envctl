package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func tagTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{Name: "alpha", Tags: []string{"prod"}},
			{Name: "beta", Tags: []string{"staging", "qa"}},
			{Name: "gamma", Tags: []string{}},
		},
	}
}

func TestTag_AddNew(t *testing.T) {
	cfg := tagTestConfig()
	if err := Tag(cfg, "gamma", []string{"dev"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("gamma")
	if len(set.Tags) != 1 || set.Tags[0] != "dev" {
		t.Errorf("expected [dev], got %v", set.Tags)
	}
}

func TestTag_DeduplicatesTags(t *testing.T) {
	cfg := tagTestConfig()
	if err := Tag(cfg, "alpha", []string{"prod", "extra"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("alpha")
	if len(set.Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", set.Tags)
	}
}

func TestTag_UnknownSet(t *testing.T) {
	cfg := tagTestConfig()
	if err := Tag(cfg, "nope", []string{"x"}); err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestTag_EmptyTagRejected(t *testing.T) {
	cfg := tagTestConfig()
	if err := Tag(cfg, "alpha", []string{""}); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestTag_NoTagsRejected(t *testing.T) {
	cfg := tagTestConfig()
	if err := Tag(cfg, "alpha", []string{}); err == nil {
		t.Fatal("expected error when no tags provided")
	}
}

func TestUntag_RemovesTag(t *testing.T) {
	cfg := tagTestConfig()
	if err := Untag(cfg, "beta", []string{"qa"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("beta")
	if len(set.Tags) != 1 || set.Tags[0] != "staging" {
		t.Errorf("expected [staging], got %v", set.Tags)
	}
}

func TestUntag_UnknownSet(t *testing.T) {
	cfg := tagTestConfig()
	if err := Untag(cfg, "nope", []string{"x"}); err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestListByTag_ReturnsMatches(t *testing.T) {
	cfg := tagTestConfig()
	names := ListByTag(cfg, "staging")
	if len(names) != 1 || names[0] != "beta" {
		t.Errorf("expected [beta], got %v", names)
	}
}

func TestListByTag_NoMatches(t *testing.T) {
	cfg := tagTestConfig()
	names := ListByTag(cfg, "unknown-tag")
	if len(names) != 0 {
		t.Errorf("expected empty, got %v", names)
	}
}
