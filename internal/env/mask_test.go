package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func maskTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"DB_PASSWORD": "secret",
					"API_KEY":     "key123",
					"APP_ENV":     "prod",
				},
				Sealed: []string{"DB_PASSWORD"},
			},
		},
	}
}

func TestMask_NewKey(t *testing.T) {
	cfg := maskTestConfig()
	res, err := Mask(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Masked) != 1 || res.Masked[0] != "API_KEY" {
		t.Errorf("expected API_KEY masked, got %v", res.Masked)
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped, got %v", res.Skipped)
	}
}

func TestMask_AlreadySealed(t *testing.T) {
	cfg := maskTestConfig()
	res, err := Mask(cfg, "production", []string{"DB_PASSWORD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Masked) != 0 {
		t.Errorf("expected nothing newly masked, got %v", res.Masked)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD skipped, got %v", res.Skipped)
	}
}

func TestMask_UnknownKey(t *testing.T) {
	cfg := maskTestConfig()
	_, err := Mask(cfg, "production", []string{"NONEXISTENT"})
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestMask_UnknownSet(t *testing.T) {
	cfg := maskTestConfig()
	_, err := Mask(cfg, "staging", []string{"API_KEY"})
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestMask_AllKeys(t *testing.T) {
	cfg := maskTestConfig()
	res, err := Mask(cfg, "production", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DB_PASSWORD already sealed → skipped; API_KEY and APP_ENV → masked
	if len(res.Masked)+len(res.Skipped) != 3 {
		t.Errorf("expected 3 total keys processed, got masked=%v skipped=%v", res.Masked, res.Skipped)
	}
}

func TestUnmask_Success(t *testing.T) {
	cfg := maskTestConfig()
	res, err := Unmask(cfg, "production", []string{"DB_PASSWORD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Masked) != 1 || res.Masked[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD unmasked, got %v", res.Masked)
	}
	set := cfg.EnvSetByName("production")
	if len(set.Sealed) != 0 {
		t.Errorf("expected sealed to be empty after unmask, got %v", set.Sealed)
	}
}

func TestUnmask_NotSealed(t *testing.T) {
	cfg := maskTestConfig()
	res, err := Unmask(cfg, "production", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "API_KEY" {
		t.Errorf("expected API_KEY skipped, got %v", res.Skipped)
	}
}

func TestMaskedValue_Sealed(t *testing.T) {
	set := &config.EnvSet{Sealed: []string{"DB_PASSWORD"}}
	got := MaskedValue(set, "DB_PASSWORD", "secret")
	if got != "******" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestMaskedValue_NotSealed(t *testing.T) {
	set := &config.EnvSet{Sealed: []string{"DB_PASSWORD"}}
	got := MaskedValue(set, "API_KEY", "key123")
	if got != "key123" {
		t.Errorf("expected plain value, got %q", got)
	}
}
