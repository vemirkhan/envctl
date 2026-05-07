package env

import (
	"fmt"
	"sort"

	"github.com/envctl/envctl/internal/config"
)

// Protect marks keys in an env set as protected, preventing them from being
// overridden by targets or modified by sync/import operations.
func Protect(cfg *config.Config, setName string, keys []string) error {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return fmt.Errorf("env set %q not found", setName)
	}

	for _, key := range keys {
		if _, ok := set.Base[key]; !ok {
			return fmt.Errorf("key %q not found in base of env set %q", key, setName)
		}
	}

	if set.Protected == nil {
		set.Protected = []string{}
	}

	existing := make(map[string]bool, len(set.Protected))
	for _, k := range set.Protected {
		existing[k] = true
	}

	for _, key := range keys {
		if !existing[key] {
			set.Protected = append(set.Protected, key)
			existing[key] = true
		}
	}

	sort.Strings(set.Protected)
	return nil
}

// Unprotect removes protection from keys in an env set.
func Unprotect(cfg *config.Config, setName string, keys []string) error {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return fmt.Errorf("env set %q not found", setName)
	}

	remove := make(map[string]bool, len(keys))
	for _, k := range keys {
		remove[k] = true
	}

	filtered := set.Protected[:0]
	for _, k := range set.Protected {
		if !remove[k] {
			filtered = append(filtered, k)
		}
	}
	set.Protected = filtered
	return nil
}

// ProtectedKeys returns the list of protected keys for a given env set.
func ProtectedKeys(cfg *config.Config, setName string) ([]string, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}
	out := make([]string, len(set.Protected))
	copy(out, set.Protected)
	return out, nil
}
