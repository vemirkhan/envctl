package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	EnvSet      string
	FromTarget  string
	ToTarget    string
	KeysPromoted []string
}

// Promote copies target-level overrides from one target to another within the
// same env set. If overwrite is false, existing keys in the destination target
// are left untouched. Returns the keys that were promoted.
func Promote(cfg *config.Config, envSetName, fromTarget, toTarget string, overwrite bool) (PromoteResult, error) {
	set := cfg.EnvSetByName(envSetName)
	if set == nil {
		return PromoteResult{}, fmt.Errorf("env set %q not found", envSetName)
	}

	srcVars, srcOK := set.Targets[fromTarget]
	if !srcOK {
		return PromoteResult{}, fmt.Errorf("source target %q not found in env set %q", fromTarget, envSetName)
	}

	if _, dstOK := set.Targets[toTarget]; !dstOK {
		return PromoteResult{}, fmt.Errorf("destination target %q not found in env set %q", toTarget, envSetName)
	}

	if set.Targets[toTarget] == nil {
		set.Targets[toTarget] = make(map[string]string)
	}

	dst := set.Targets[toTarget]
	var promoted []string

	for k, v := range srcVars {
		if _, exists := dst[k]; exists && !overwrite {
			continue
		}
		dst[k] = v
		promoted = append(promoted, k)
	}

	return PromoteResult{
		EnvSet:      envSetName,
		FromTarget:  fromTarget,
		ToTarget:    toTarget,
		KeysPromoted: promoted,
	}, nil
}
