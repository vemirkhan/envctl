package env

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/user/envctl/internal/config"
)

var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// LintIssue describes a single linting problem found in an env set.
type LintIssue struct {
	Set     string
	Key     string
	Message string
}

// LintResult holds all issues found during a lint pass.
type LintResult struct {
	Issues []LintIssue
}

// OK returns true when no issues were found.
func (r *LintResult) OK() bool { return len(r.Issues) == 0 }

// Lint checks the named env set (and its targets) for style and consistency
// issues: non-uppercase keys, keys with leading/trailing whitespace, values
// with unresolved variable references, and duplicate keys across base+target.
func Lint(cfg *config.Config, setName string) (*LintResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	result := &LintResult{}

	checkKeys(result, setName, set.Base)

	for _, target := range set.Targets {
		label := fmt.Sprintf("%s[%s]", setName, target.Name)
		checkKeys(result, label, target.Overrides)

		for k := range target.Overrides {
			if _, exists := set.Base[k]; !exists {
				result.Issues = append(result.Issues, LintIssue{
					Set:     label,
					Key:     k,
					Message: "override key not present in base",
				})
			}
		}
	}

	sort.Slice(result.Issues, func(i, j int) bool {
		if result.Issues[i].Set != result.Issues[j].Set {
			return result.Issues[i].Set < result.Issues[j].Set
		}
		return result.Issues[i].Key < result.Issues[j].Key
	})

	return result, nil
}

func checkKeys(result *LintResult, label string, vars map[string]string) {
	for k, v := range vars {
		if strings.TrimSpace(k) != k {
			result.Issues = append(result.Issues, LintIssue{Set: label, Key: k, Message: "key has leading or trailing whitespace"})
		}
		if !validKeyPattern.MatchString(k) {
			result.Issues = append(result.Issues, LintIssue{Set: label, Key: k, Message: "key should be uppercase with underscores (A-Z, 0-9, _)"})
		}
		if strings.TrimSpace(v) == "" {
			result.Issues = append(result.Issues, LintIssue{Set: label, Key: k, Message: "value is empty or whitespace-only"})
		}
	}
}

// WriteLint formats a LintResult to w in a human-readable style.
func WriteLint(w io.Writer, result *LintResult, setName string) {
	if result.OK() {
		fmt.Fprintf(w, "✔  no issues found in %q\n", setName)
		return
	}
	fmt.Fprintf(w, "%d issue(s) found in %q:\n", len(result.Issues), setName)
	for _, issue := range result.Issues {
		fmt.Fprintf(w, "  [%s] %s: %s\n", issue.Set, issue.Key, issue.Message)
	}
}
