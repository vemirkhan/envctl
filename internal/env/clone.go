package env

import (
	"fmt"

	"github.com/envctl/envctl/internal/config"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	Source      string
	Destination string
	KeysCopied  int
}

// Clone creates a deep copy of an existing env set under a new name,
// including all base vars and optionally all target overrides.
// If includeTargets is false, only the base vars are cloned.
func Clone(cfg *config.Config, srcName, dstName string, includeTargets bool) (CloneResult, error) {
	src := cfg.EnvSetByName(srcName)
	if src == nil {
		return CloneResult{}, fmt.Errorf("env set %q not found", srcName)
	}

	if cfg.EnvSetByName(dstName) != nil {
		return CloneResult{}, fmt.Errorf("env set %q already exists", dstName)
	}

	newSet := config.EnvSet{
		Name: dstName,
		Base: make(map[string]string, len(src.Base)),
	}

	for k, v := range src.Base {
		newSet.Base[k] = v
	}

	if includeTargets && len(src.Targets) > 0 {
		newSet.Targets = make([]config.Target, len(src.Targets))
		for i, t := range src.Targets {
			clonedTarget := config.Target{
				Name:      t.Name,
				Overrides: make(map[string]string, len(t.Overrides)),
			}
			for k, v := range t.Overrides {
				clonedTarget.Overrides[k] = v
			}
			newSet.Targets[i] = clonedTarget
		}
	}

	cfg.EnvSets = append(cfg.EnvSets, newSet)

	return CloneResult{
		Source:      srcName,
		Destination: dstName,
		KeysCopied:  len(newSet.Base),
	}, nil
}
