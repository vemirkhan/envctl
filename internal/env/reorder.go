package env

import (
	"fmt"

	"github.com/envctl/envctl/internal/config"
)

// Reorder changes the position of keys within an env set's base map by
// rebuilding the ordered representation. Since Go maps are unordered, reorder
// updates a dedicated KeyOrder slice on the EnvSet so exporters can respect it.
func Reorder(cfg *config.Config, setName string, orderedKeys []string) error {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return fmt.Errorf("env set %q not found", setName)
	}

	// Validate all provided keys exist in base.
	for _, k := range orderedKeys {
		if _, ok := set.Base[k]; !ok {
			return fmt.Errorf("key %q not found in base of env set %q", k, setName)
		}
	}

	// Ensure no duplicates in the provided order.
	seen := make(map[string]bool, len(orderedKeys))
	for _, k := range orderedKeys {
		if seen[k] {
			return fmt.Errorf("duplicate key %q in reorder list", k)
		}
		seen[k] = true
	}

	// Append any base keys not explicitly listed at the end (stable).
	existing := set.KeyOrder
	if len(existing) == 0 {
		for k := range set.Base {
			existing = append(existing, k)
		}
	}

	appended := []string{}
	for _, k := range existing {
		if !seen[k] {
			appended = append(appended, k)
		}
	}

	set.KeyOrder = append(orderedKeys, appended...)
	return nil
}
