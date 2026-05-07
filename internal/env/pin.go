package env

import (
	"fmt"

	"github.com/envctl/envctl/internal/config"
)

// PinResult holds the outcome of a pin or unpin operation.
type PinResult struct {
	Set     string
	Key     string
	Pinned  bool
	Skipped bool
	Reason  string
}

// Pin marks one or more keys in an env set as pinned, preventing them from
// being overwritten by sync or promote operations.
func Pin(cfg *config.Config, setName string, keys []string) ([]PinResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	if set.Pinned == nil {
		set.Pinned = []string{}
	}

	pinned := make(map[string]struct{}, len(set.Pinned))
	for _, k := range set.Pinned {
		pinned[k] = struct{}{}
	}

	var results []PinResult
	for _, key := range keys {
		if _, exists := set.Base[key]; !exists {
			results = append(results, PinResult{Set: setName, Key: key, Skipped: true, Reason: "key not in base"})
			continue
		}
		if _, already := pinned[key]; already {
			results = append(results, PinResult{Set: setName, Key: key, Skipped: true, Reason: "already pinned"})
			continue
		}
		set.Pinned = append(set.Pinned, key)
		pinned[key] = struct{}{}
		results = append(results, PinResult{Set: setName, Key: key, Pinned: true})
	}

	return results, nil
}

// Unpin removes pin marks from one or more keys in an env set.
func Unpin(cfg *config.Config, setName string, keys []string) ([]PinResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	remove := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		remove[k] = struct{}{}
	}

	var next []string
	unpinned := make(map[string]struct{})
	for _, k := range set.Pinned {
		if _, ok := remove[k]; ok {
			unpinned[k] = struct{}{}
		} else {
			next = append(next, k)
		}
	}
	set.Pinned = next

	var results []PinResult
	for _, key := range keys {
		if _, was := unpinned[key]; was {
			results = append(results, PinResult{Set: setName, Key: key, Pinned: false})
		} else {
			results = append(results, PinResult{Set: setName, Key: key, Skipped: true, Reason: "not pinned"})
		}
	}
	return results, nil
}
