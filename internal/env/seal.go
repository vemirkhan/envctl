package env

import (
	"fmt"
	"sort"

	"github.com/envctl/envctl/internal/config"
)

// SealResult holds the outcome of a seal operation.
type SealResult struct {
	Set    string
	Sealed []string
}

// Seal marks the given keys in an env set as sealed (read-only), preventing
// them from being overridden by any target. If keys is empty, all base keys
// are sealed.
func Seal(cfg *config.Config, setName string, keys []string) (*SealResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	if set.Sealed == nil {
		set.Sealed = []string{}
	}

	targetKeys := keys
	if len(targetKeys) == 0 {
		for k := range set.Base {
			targetKeys = append(targetKeys, k)
		}
		sort.Strings(targetKeys)
	}

	sealed := []string{}
	existing := make(map[string]bool, len(set.Sealed))
	for _, k := range set.Sealed {
		existing[k] = true
	}

	for _, k := range targetKeys {
		if _, ok := set.Base[k]; !ok {
			return nil, fmt.Errorf("key %q not found in base of env set %q", k, setName)
		}
		if !existing[k] {
			set.Sealed = append(set.Sealed, k)
			sealed = append(sealed, k)
			existing[k] = true
		}
	}
	sort.Strings(set.Sealed)

	return &SealResult{Set: setName, Sealed: sealed}, nil
}

// Unseal removes the sealed designation from the given keys. If keys is empty,
// all sealed keys are removed.
func Unseal(cfg *config.Config, setName string, keys []string) (*SealResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	remove := make(map[string]bool)
	if len(keys) == 0 {
		for _, k := range set.Sealed {
			remove[k] = true
		}
	} else {
		for _, k := range keys {
			remove[k] = true
		}
	}

	unsealedKeys := []string{}
	for k := range remove {
		unsealedKeys = append(unsealedKeys, k)
	}
	sort.Strings(unsealedKeys)

	updated := []string{}
	for _, k := range set.Sealed {
		if !remove[k] {
			updated = append(updated, k)
		}
	}
	set.Sealed = updated

	return &SealResult{Set: setName, Sealed: unsealedKeys}, nil
}
