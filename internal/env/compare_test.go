package env_test

import (
	"bytes"
	"testing"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
)

func compareTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{"HOST": "localhost", "PORT": "8080", "DEBUG": "true"},
				Targets: map[string]map[string]string{
					"prod": {"HOST": "prod.example.com", "DEBUG": "false"},
				},
			},
			{
				Name: "worker",
				Base: map[string]string{"HOST": "localhost", "PORT": "9090", "QUEUE": "default"},
			},
		},
	}
}

func TestCompare_SameSets(t *testing.T) {
	cfg := compareTestConfig()
	r, err := env.Compare(cfg, "app", "app", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Differ) != 0 || len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Errorf("expected identical sets, got diffs")
	}
}

func TestCompare_DifferentSets(t *testing.T) {
	cfg := compareTestConfig()
	r, err := env.Compare(cfg, "app", "worker", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Differ["PORT"]; !ok {
		t.Errorf("expected PORT to differ")
	}
	if _, ok := r.OnlyInA["DEBUG"]; !ok {
		t.Errorf("expected DEBUG only in app")
	}
	if _, ok := r.OnlyInB["QUEUE"]; !ok {
		t.Errorf("expected QUEUE only in worker")
	}
	if _, ok := r.Same["HOST"]; !ok {
		t.Errorf("expected HOST to be same")
	}
}

func TestCompare_WithTarget(t *testing.T) {
	cfg := compareTestConfig()
	r, err := env.Compare(cfg, "app", "worker", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Differ["HOST"][0] != "prod.example.com" {
		t.Errorf("expected prod HOST for app, got %q", r.Differ["HOST"][0])
	}
}

func TestCompare_UnknownSet(t *testing.T) {
	cfg := compareTestConfig()
	_, err := env.Compare(cfg, "app", "ghost", "")
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestWriteCompare_Output(t *testing.T) {
	cfg := compareTestConfig()
	r, _ := env.Compare(cfg, "app", "worker", "")
	var buf bytes.Buffer
	env.WriteCompare(&buf, r)
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("Comparing")) {
		t.Errorf("expected header in output")
	}
	if !bytes.Contains([]byte(out), []byte("PORT")) {
		t.Errorf("expected PORT in output")
	}
}

func TestWriteCompare_NoDiff(t *testing.T) {
	cfg := compareTestConfig()
	r, _ := env.Compare(cfg, "app", "app", "")
	var buf bytes.Buffer
	env.WriteCompare(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("No differences")) {
		t.Errorf("expected no-differences message")
	}
}
