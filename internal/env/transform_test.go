package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func transformTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "app",
				Base: map[string]string{
					"APP_NAME": "myapp",
					"LOG_LEVEL": "  debug  ",
					"REGION": "us-east-1",
				},
				Targets: map[string]map[string]string{
					"prod": {
						"LOG_LEVEL": "warn",
					},
				},
			},
		},
	}
}

func TestTransform_UpperBase(t *testing.T) {
	cfg := transformTestConfig()
	changed, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Keys:    []string{"APP_NAME"},
		Op:      TransformUpper,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := changed["APP_NAME"]; ok {
		t.Error("APP_NAME already uppercase, should not appear in changed")
	}
}

func TestTransform_LowerBase(t *testing.T) {
	cfg := transformTestConfig()
	changed, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Keys:    []string{"REGION"},
		Op:      TransformLower,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed["REGION"] != "us-east-1" {
		t.Errorf("expected us-east-1, got %q", changed["REGION"])
	}
}

func TestTransform_TrimWS(t *testing.T) {
	cfg := transformTestConfig()
	changed, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Keys:    []string{"LOG_LEVEL"},
		Op:      TransformTrimWS,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed["LOG_LEVEL"] != "debug" {
		t.Errorf("expected 'debug', got %q", changed["LOG_LEVEL"])
	}
}

func TestTransform_DryRun(t *testing.T) {
	cfg := transformTestConfig()
	_, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Keys:    []string{"LOG_LEVEL"},
		Op:      TransformTrimWS,
		DryRun:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set := cfg.EnvSetByName("app")
	if set.Base["LOG_LEVEL"] != "  debug  " {
		t.Error("dry-run should not modify the config")
	}
}

func TestTransform_WithTarget(t *testing.T) {
	cfg := transformTestConfig()
	changed, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Target:  "prod",
		Keys:    []string{"LOG_LEVEL"},
		Op:      TransformUpper,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed["LOG_LEVEL"] != "WARN" {
		t.Errorf("expected WARN, got %q", changed["LOG_LEVEL"])
	}
}

func TestTransform_UnknownSet(t *testing.T) {
	cfg := transformTestConfig()
	_, err := Transform(cfg, TransformOptions{SetName: "nope", Op: TransformUpper})
	if err == nil {
		t.Error("expected error for unknown set")
	}
}

func TestTransform_InvalidOp(t *testing.T) {
	cfg := transformTestConfig()
	_, err := Transform(cfg, TransformOptions{SetName: "app", Op: "rot13"})
	if err == nil {
		t.Error("expected error for invalid op")
	}
}

func TestTransform_UnknownKey(t *testing.T) {
	cfg := transformTestConfig()
	_, err := Transform(cfg, TransformOptions{
		SetName: "app",
		Keys:    []string{"MISSING"},
		Op:      TransformUpper,
	})
	if err == nil {
		t.Error("expected error for unknown key")
	}
}
