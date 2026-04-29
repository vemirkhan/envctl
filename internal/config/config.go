package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Target describes a deployment target within an env set.
type Target struct {
	Name      string            `yaml:"name"`
	File      string            `yaml:"file"`
	Format    string            `yaml:"format"`
	Overrides map[string]string `yaml:"overrides"`
}

// EnvSet represents a named collection of environment variables with optional
// per-target overrides.
type EnvSet struct {
	Name    string            `yaml:"name"`
	Vars    map[string]string `yaml:"vars"`
	Targets []Target          `yaml:"targets"`
}

// Config is the top-level configuration structure.
type Config struct {
	EnvSets []EnvSet `yaml:"env_sets"`
}

// Load reads and parses the YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	seen := make(map[string]bool)
	for _, es := range c.EnvSets {
		if seen[es.Name] {
			return fmt.Errorf("duplicate env set name: %q", es.Name)
		}
		seen[es.Name] = true
	}
	return nil
}

// EnvSetByName returns the EnvSet with the given name or an error.
func (c *Config) EnvSetByName(name string) (*EnvSet, error) {
	for i := range c.EnvSets {
		if c.EnvSets[i].Name == name {
			return &c.EnvSets[i], nil
		}
	}
	return nil, fmt.Errorf("env set %q not found", name)
}
