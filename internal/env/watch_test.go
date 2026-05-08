package env

import (
	"testing"
	"time"

	"github.com/envctl/envctl/internal/config"
)

func watchTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"HOST": "localhost",
					"PORT": "8080",
				},
				Targets: map[string]map[string]string{
					"prod": {"HOST": "prod.example.com"},
				},
			},
		},
	}
}

func TestWatch_NoChange(t *testing.T) {
	cfg := watchTestConfig()
	done := make(chan struct{})

	opts := WatchOptions{Set: "app", Interval: 10 * time.Millisecond, MaxPolls: 3}
	results, errs := Watch(cfg, opts, done)

	var got []WatchResult
	for r := range results {
		got = append(got, r)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no changes, got %d", len(got))
	}
}

func TestWatch_UnknownSet(t *testing.T) {
	cfg := watchTestConfig()
	done := make(chan struct{})

	opts := WatchOptions{Set: "nonexistent", Interval: 10 * time.Millisecond, MaxPolls: 1}
	_, errs := Watch(cfg, opts, done)

	err := <-errs
	if err == nil {
		t.Fatal("expected error for unknown set, got nil")
	}
}

func TestWatch_DoneSignal(t *testing.T) {
	cfg := watchTestConfig()
	done := make(chan struct{})

	opts := WatchOptions{Set: "app", Interval: 5 * time.Millisecond, MaxPolls: 0}
	results, errs := Watch(cfg, opts, done)

	time.AfterFunc(25*time.Millisecond, func() { close(done) })

	for range results {
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatch_WithTarget(t *testing.T) {
	cfg := watchTestConfig()
	done := make(chan struct{})

	opts := WatchOptions{Set: "app", Target: "prod", Interval: 10 * time.Millisecond, MaxPolls: 2}
	results, errs := Watch(cfg, opts, done)

	for range results {
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error watching with target: %v", err)
	}
}
