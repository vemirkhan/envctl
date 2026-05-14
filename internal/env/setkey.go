package env

import (
	"fmt"

	"github.com/envctl/envctl/internal/config"
)

// SetKeyOptions controls behaviour of SetKey.
type SetKeyOptions struct {
	Target    string
	Overwrite bool
}

// SetKey sets or updates a single key in an env set's base vars or a specific
// target override. Returns an error if the key already exists and Overwrite is
// false, or if the named set / target cannot be found.
func SetKey(cfg *config.Config, setName, key, value string, opts SetKeyOptions) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}

	es := cfg.EnvSetByName(setName)
	if es == nil {
		return fmt.Errorf("env set %q not found", setName)
	}

	if opts.Target == "" {
		if es.Base == nil {
			es.Base = make(map[string]string)
		}
		if _, exists := es.Base[key]; exists && !opts.Overwrite {
			return fmt.Errorf("key %q already exists in base of %q; use --overwrite to replace", key, setName)
		}
		es.Base[key] = value
		return nil
	}

	// Target override path.
	for i := range es.Targets {
		if es.Targets[i].Name == opts.Target {
			if es.Targets[i].Overrides == nil {
				es.Targets[i].Overrides = make(map[string]string)
			}
			if _, exists := es.Targets[i].Overrides[key]; exists && !opts.Overwrite {
				return fmt.Errorf("key %q already exists in target %q of %q; use --overwrite to replace", key, opts.Target, setName)
			}
			es.Targets[i].Overrides[key] = value
			return nil
		}
	}
	return fmt.Errorf("target %q not found in env set %q", opts.Target, setName)
}
