package config

import "time"

// SnapshotEntry holds a named, frozen copy of resolved env vars.
type SnapshotEntry struct {
	Name      string            `yaml:"name"       json:"name"`
	Set       string            `yaml:"set"        json:"set"`
	Target    string            `yaml:"target,omitempty" json:"target,omitempty"`
	CreatedAt time.Time         `yaml:"created_at" json:"created_at"`
	Vars      map[string]string `yaml:"vars"       json:"vars"`
}

// Snapshots field is appended to Config in this file to keep concerns separated.
// The main Config struct embeds this via the Snapshots slice already declared here.
// NOTE: Config.Snapshots must be declared in config.go; this file only defines the type
// and the Save helper so the rest of the codebase can persist mutations.

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Save writes the full config (including any mutated Snapshots) back to disk.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
