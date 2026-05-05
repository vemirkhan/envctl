package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// Delete removes an env set by name from the config.
// If removeRefs is true, any target overrides referencing the set are also removed.
// Returns an error if the set does not exist.
func Delete(cfg *config.Config, name string, removeRefs bool) error {
	idx := -1
	for i, s := range cfg.EnvSets {
		if s.Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("env set %q not found", name)
	}

	// Remove the env set.
	cfg.EnvSets = append(cfg.EnvSets[:idx], cfg.EnvSets[idx+1:]...)

	if !removeRefs {
		return nil
	}

	// Remove any target override entries that reference the deleted set.
	for ti := range cfg.Targets {
		filtered := cfg.Targets[ti].Overrides[:0]
		for _, ov := range cfg.Targets[ti].Overrides {
			if ov.EnvSet != name {
				filtered = append(filtered, ov)
			}
		}
		cfg.Targets[ti].Overrides = filtered
	}

	return nil
}
