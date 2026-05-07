package env

import (
	"fmt"
	"sort"

	"github.com/envctl/envctl/internal/config"
)

// ExtractResult holds the result of an extract operation.
type ExtractResult struct {
	SetName string
	Keys    []string
	Vars    map[string]string
}

// Extract pulls a subset of keys from an env set (and optional target) into a
// new standalone env set. If destName already exists, an error is returned
// unless overwrite is true.
func Extract(cfg *config.Config, srcName string, keys []string, destName string, target string, overwrite bool) (*ExtractResult, error) {
	src := cfg.EnvSetByName(srcName)
	if src == nil {
		return nil, fmt.Errorf("env set %q not found", srcName)
	}

	if !overwrite && cfg.EnvSetByName(destName) != nil {
		return nil, fmt.Errorf("env set %q already exists; use --overwrite to replace", destName)
	}

	resolved, err := Resolve(cfg, srcName, target)
	if err != nil {
		return nil, err
	}

	// Validate all requested keys exist in the resolved set.
	for _, k := range keys {
		if _, ok := resolved[k]; !ok {
			return nil, fmt.Errorf("key %q not found in env set %q (target: %q)", k, srcName, target)
		}
	}

	vars := make(map[string]string, len(keys))
	for _, k := range keys {
		vars[k] = resolved[k]
	}

	// Remove existing destination if overwrite.
	for i, es := range cfg.EnvSets {
		if es.Name == destName {
			cfg.EnvSets = append(cfg.EnvSets[:i], cfg.EnvSets[i+1:]...)
			break
		}
	}

	cfg.EnvSets = append(cfg.EnvSets, config.EnvSet{
		Name: destName,
		Base: vars,
	})

	sortedKeys := make([]string, 0, len(vars))
	for k := range vars {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	return &ExtractResult{
		SetName: destName,
		Keys:    sortedKeys,
		Vars:    vars,
	}, nil
}
