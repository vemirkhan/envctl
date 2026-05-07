package env

import (
	"fmt"
	"strings"

	"github.com/your-org/envctl/internal/config"
)

// TrimResult holds the result of a trim operation.
type TrimResult struct {
	Set     string
	Removed []string
}

// Trim removes all keys from an env set's base (and optionally target overrides)
// that match the given prefix or suffix pattern.
func Trim(cfg *config.Config, setName, prefix, suffix string, targets []string) (TrimResult, error) {
	if prefix == "" && suffix == "" {
		return TrimResult{}, fmt.Errorf("at least one of --prefix or --suffix must be specified")
	}

	set := cfg.EnvSetByName(setName)
	if set == nil {
		return TrimResult{}, fmt.Errorf("env set %q not found", setName)
	}

	var removed []string

	// Trim base keys
	for k := range set.Base {
		if matchesTrim(k, prefix, suffix) {
			delete(set.Base, k)
			removed = append(removed, k)
		}
	}

	// Trim target overrides
	for _, tname := range targets {
		for _, t := range set.Targets {
			if t.Name != tname {
				continue
			}
			for k := range t.Overrides {
				if matchesTrim(k, prefix, suffix) {
					delete(t.Overrides, k)
					removed = append(removed, fmt.Sprintf("%s[%s]", k, tname))
				}
			}
		}
	}

	return TrimResult{Set: setName, Removed: removed}, nil
}

func matchesTrim(key, prefix, suffix string) bool {
	if prefix != "" && strings.HasPrefix(key, prefix) {
		return true
	}
	if suffix != "" && strings.HasSuffix(key, suffix) {
		return true
	}
	return false
}
