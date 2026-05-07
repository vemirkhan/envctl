package env

import (
	"testing"

	"github.com/envctl/envctl/internal/config"
)

func extractTestConfig() *config.Config {
	return &config.Config{
		EnvSets: []config.EnvSet{
			{
				Name: "production",
				Base: map[string]string{
					"DB_HOST": "prod-db",
					"DB_PORT": "5432",
					"API_KEY": "secret",
					"LOG_LEVEL": "warn",
				},
				Targets: map[string]map[string]string{
					"eu": {"DB_HOST": "eu-db"},
				},
			},
		},
	}
}

func TestExtract_BasicKeys(t *testing.T) {
	cfg := extractTestConfig()
	res, err := Extract(cfg, "production", []string{"DB_HOST", "DB_PORT"}, "db-only", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SetName != "db-only" {
		t.Errorf("expected set name db-only, got %s", res.SetName)
	}
	if len(res.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(res.Vars))
	}
	if res.Vars["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %s", res.Vars["DB_HOST"])
	}
	if cfg.EnvSetByName("db-only") == nil {
		t.Error("expected db-only set to be added to config")
	}
}

func TestExtract_WithTarget(t *testing.T) {
	cfg := extractTestConfig()
	res, err := Extract(cfg, "production", []string{"DB_HOST"}, "eu-db", "eu", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["DB_HOST"] != "eu-db" {
		t.Errorf("expected eu-db, got %s", res.Vars["DB_HOST"])
	}
}

func TestExtract_UnknownSet(t *testing.T) {
	cfg := extractTestConfig()
	_, err := Extract(cfg, "staging", []string{"DB_HOST"}, "out", "", false)
	if err == nil {
		t.Fatal("expected error for unknown set")
	}
}

func TestExtract_UnknownKey(t *testing.T) {
	cfg := extractTestConfig()
	_, err := Extract(cfg, "production", []string{"MISSING_KEY"}, "out", "", false)
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestExtract_DestinationExists(t *testing.T) {
	cfg := extractTestConfig()
	cfg.EnvSets = append(cfg.EnvSets, config.EnvSet{Name: "existing", Base: map[string]string{}})
	_, err := Extract(cfg, "production", []string{"DB_HOST"}, "existing", "", false)
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestExtract_Overwrite(t *testing.T) {
	cfg := extractTestConfig()
	cfg.EnvSets = append(cfg.EnvSets, config.EnvSet{Name: "existing", Base: map[string]string{"OLD": "val"}})
	res, err := Extract(cfg, "production", []string{"API_KEY"}, "existing", "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["OLD"]; ok {
		t.Error("expected OLD key to be gone after overwrite")
	}
	if res.Vars["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %s", res.Vars["API_KEY"])
	}
}
