package env

import (
	"fmt"
	"sort"

	"github.com/envctl/envctl/internal/config"
)

// TagResult holds the result of a tag operation.
type TagResult struct {
	Set  string
	Tags []string
}

// Tag adds one or more tags to an env set. Tags are stored as metadata
// on the EnvSet and can be used for filtering and grouping.
func Tag(cfg *config.Config, setName string, tags []string) error {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return fmt.Errorf("env set %q not found", setName)
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag must be provided")
	}
	existing := make(map[string]struct{}, len(set.Tags))
	for _, t := range set.Tags {
		existing[t] = struct{}{}
	}
	for _, t := range tags {
		if t == "" {
			return fmt.Errorf("tag must not be empty")
		}
		existing[t] = struct{}{}
	}
	merged := make([]string, 0, len(existing))
	for t := range existing {
		merged = append(merged, t)
	}
	sort.Strings(merged)
	set.Tags = merged
	return nil
}

// Untag removes one or more tags from an env set.
func Untag(cfg *config.Config, setName string, tags []string) error {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return fmt.Errorf("env set %q not found", setName)
	}
	remove := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		remove[t] = struct{}{}
	}
	filtered := set.Tags[:0]
	for _, t := range set.Tags {
		if _, skip := remove[t]; !skip {
			filtered = append(filtered, t)
		}
	}
	set.Tags = filtered
	return nil
}

// ListByTag returns all env set names that carry the given tag.
func ListByTag(cfg *config.Config, tag string) []string {
	var names []string
	for i := range cfg.EnvSets {
		for _, t := range cfg.EnvSets[i].Tags {
			if t == tag {
				names = append(names, cfg.EnvSets[i].Name)
				break
			}
		}
	}
	sort.Strings(names)
	return names
}
