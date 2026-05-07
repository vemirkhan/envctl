package env

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
)

// MaskResult holds the result of a mask operation for a single env set.
type MaskResult struct {
	SetName string
	Masked  []string
	Skipped []string
}

// Mask marks the given keys as sealed (sensitive) in the specified env set.
// If keys is empty, all keys in the base are masked.
// Returns an error if the set is not found or a key does not exist in base.
func Mask(cfg *config.Config, setName string, keys []string) (*MaskResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	result := &MaskResult{SetName: setName}

	if len(keys) == 0 {
		for k := range set.Base {
			keys = append(keys, k)
		}
	}

	existing := make(map[string]bool, len(set.Sealed))
	for _, k := range set.Sealed {
		existing[k] = true
	}

	for _, k := range keys {
		if _, ok := set.Base[k]; !ok {
			return nil, fmt.Errorf("key %q not found in base of set %q", k, setName)
		}
		if existing[k] {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		set.Sealed = append(set.Sealed, k)
		existing[k] = true
		result.Masked = append(result.Masked, k)
	}

	return result, nil
}

// Unmask removes the given keys from the sealed list of the specified env set.
func Unmask(cfg *config.Config, setName string, keys []string) (*MaskResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	result := &MaskResult{SetName: setName}

	remove := make(map[string]bool, len(keys))
	for _, k := range keys {
		remove[k] = true
	}

	var remaining []string
	for _, k := range set.Sealed {
		if remove[k] {
			result.Masked = append(result.Masked, k)
		} else {
			remaining = append(remaining, k)
		}
	}

	for k := range remove {
		found := false
		for _, mk := range result.Masked {
			if mk == k {
				found = true
				break
			}
		}
		if !found {
			result.Skipped = append(result.Skipped, k)
		}
	}

	set.Sealed = remaining
	return result, nil
}

// MaskedValue returns the display value for a key, masking it if sealed.
func MaskedValue(set *config.EnvSet, key, value string) string {
	for _, k := range set.Sealed {
		if k == key {
			return strings.Repeat("*", len(value))
		}
	}
	return value
}
