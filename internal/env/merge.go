package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// MergeResult holds the merged environment variables and metadata.
type MergeResult struct {
	Vars     map[string]string
	Sources  map[string]string // key -> which set it came from
	Conflicts []MergeConflict
}

// MergeConflict describes a key that existed in multiple sets with different values.
type MergeConflict struct {
	Key    string
	Values map[string]string // setName -> value
}

// Merge combines multiple env sets into a single variable map.
// Later sets take precedence over earlier ones. Conflicts are recorded
// but do not cause an error; the last set wins.
func Merge(cfg *config.Config, setNames []string) (*MergeResult, error) {
	if len(setNames) == 0 {
		return nil, fmt.Errorf("merge: at least one env set name required")
	}

	result := &MergeResult{
		Vars:    make(map[string]string),
		Sources: make(map[string]string),
	}

	// track per-key values across sets for conflict detection
	seen := make(map[string]map[string]string) // key -> setName -> value

	for _, name := range setNames {
		set := cfg.EnvSetByName(name)
		if set == nil {
			return nil, fmt.Errorf("merge: unknown env set %q", name)
		}

		for k, v := range set.Base {
			if _, exists := seen[k]; !exists {
				seen[k] = make(map[string]string)
			}
			seen[k][name] = v
			result.Vars[k] = v
			result.Sources[k] = name
		}
	}

	for k, bySet := range seen {
		if len(bySet) > 1 {
			// check if values actually differ
			unique := make(map[string]struct{})
			for _, v := range bySet {
				unique[v] = struct{}{}
			}
			if len(unique) > 1 {
				result.Conflicts = append(result.Conflicts, MergeConflict{
					Key:    k,
					Values: bySet,
				})
			}
		}
	}

	return result, nil
}
