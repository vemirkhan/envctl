package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvSet represents a named set of environment variables.
type EnvSet struct {
	Name      string            `yaml:"name"`
	Variables map[string]string `yaml:"variables"`
}

// Target represents a deployment target (e.g., staging, production).
type Target struct {
	Name   string   `yaml:"name"`
	Envs   []string `yaml:"envs"`
}

// Config is the top-level envctl configuration structure.
type Config struct {
	Version string    `yaml:"version"`
	EnvSets []EnvSet  `yaml:"env_sets"`
	Targets []Target  `yaml:"targets"`
}

// Load reads and parses an envctl config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config to the given path in YAML format.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// EnvSetByName returns the EnvSet with the given name, or an error if not found.
func (c *Config) EnvSetByName(name string) (*EnvSet, error) {
	for i := range c.EnvSets {
		if c.EnvSets[i].Name == name {
			return &c.EnvSets[i], nil
		}
	}
	return nil, fmt.Errorf("env set %q not found", name)
}

func (c *Config) validate() error {
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}
	seen := make(map[string]bool)
	for _, es := range c.EnvSets {
		if es.Name == "" {
			return fmt.Errorf("env set name must not be empty")
		}
		if seen[es.Name] {
			return fmt.Errorf("duplicate env set name %q", es.Name)
		}
		seen[es.Name] = true
	}
	return nil
}
