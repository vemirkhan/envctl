package env

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/user/envctl/internal/config"
)

// GrepResult holds a single match found during a grep operation.
type GrepResult struct {
	SetName string
	Target  string // empty string means base
	Key     string
	Value   string
}

// Grep searches for keys or values matching pattern across one or all env sets.
// If setName is empty, all sets are searched. matchValues controls whether
// values are also searched in addition to keys.
func Grep(cfg *config.Config, setName, pattern string, matchValues bool) ([]GrepResult, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}

	var sets []*config.EnvSet
	if setName != "" {
		s := cfg.EnvSetByName(setName)
		if s == nil {
			return nil, fmt.Errorf("env set %q not found", setName)
		}
		sets = append(sets, s)
	} else {
		sets = cfg.EnvSets
	}

	var results []GrepResult
	for _, s := range sets {
		for k, v := range s.Base {
			if re.MatchString(k) || (matchValues && re.MatchString(v)) {
				results = append(results, GrepResult{SetName: s.Name, Target: "", Key: k, Value: v})
			}
		}
		for _, t := range s.Targets {
			for k, v := range t.Overrides {
				if re.MatchString(k) || (matchValues && re.MatchString(v)) {
					results = append(results, GrepResult{SetName: s.Name, Target: t.Name, Key: k, Value: v})
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].SetName != results[j].SetName {
			return results[i].SetName < results[j].SetName
		}
		if results[i].Target != results[j].Target {
			return results[i].Target < results[j].Target
		}
		return results[i].Key < results[j].Key
	})
	return results, nil
}

// WriteGrep writes grep results in a human-readable format to w.
func WriteGrep(w io.Writer, results []GrepResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no matches found")
		return
	}
	for _, r := range results {
		location := r.SetName
		if r.Target != "" {
			location = r.SetName + ":" + r.Target
		}
		fmt.Fprintf(w, "[%s] %s=%s\n", location, r.Key, strings.TrimSpace(r.Value))
	}
}
