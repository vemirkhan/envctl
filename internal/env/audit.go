package env

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envctl/envctl/internal/config"
)

// AuditResult holds findings for a single env set.
type AuditResult struct {
	SetName        string
	UnusedOverrides []string // target keys that shadow base with identical value
	MissingInBase   []string // target keys not present in base
	EmptyValues     []string // keys with empty values (base or target)
}

// Audit inspects an env set for common issues such as redundant target
// overrides, keys missing from base, and empty values.
func Audit(cfg *config.Config, setName string) ([]AuditResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	var results []AuditResult

	// Collect empty values in base.
	var emptyBase []string
	for k, v := range set.Base {
		if strings.TrimSpace(v) == "" {
			emptyBase = append(emptyBase, k)
		}
	}
	sort.Strings(emptyBase)

	baseResult := AuditResult{
		SetName:     setName + " (base)",
		EmptyValues: emptyBase,
	}
	results = append(results, baseResult)

	for _, target := range set.Targets {
		var unused, missing, empty []string

		for k, v := range target.Overrides {
			baseVal, inBase := set.Base[k]
			if !inBase {
				missing = append(missing, k)
			} else if baseVal == v {
				unused = append(unused, k)
			}
			if strings.TrimSpace(v) == "" {
				empty = append(empty, k)
			}
		}

		sort.Strings(unused)
		sort.Strings(missing)
		sort.Strings(empty)

		results = append(results, AuditResult{
			SetName:         fmt.Sprintf("%s (target: %s)", setName, target.Name),
			UnusedOverrides: unused,
			MissingInBase:   missing,
			EmptyValues:     empty,
		})
	}

	return results, nil
}

// WriteAudit writes human-readable audit output to w.
func WriteAudit(w io.Writer, results []AuditResult) {
	for _, r := range results {
		hasIssues := len(r.UnusedOverrides) > 0 || len(r.MissingInBase) > 0 || len(r.EmptyValues) > 0
		if !hasIssues {
			fmt.Fprintf(w, "[%s] OK\n", r.SetName)
			continue
		}
		fmt.Fprintf(w, "[%s]\n", r.SetName)
		for _, k := range r.UnusedOverrides {
			fmt.Fprintf(w, "  REDUNDANT override: %s (same as base)\n", k)
		}
		for _, k := range r.MissingInBase {
			fmt.Fprintf(w, "  UNKNOWN key: %s (not in base)\n", k)
		}
		for _, k := range r.EmptyValues {
			fmt.Fprintf(w, "  EMPTY value: %s\n", k)
		}
	}
}
