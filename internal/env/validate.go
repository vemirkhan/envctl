package env

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envctl/internal/config"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds a list of validation issues found in an env set.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d issue(s):\n  %s",
		len(e.Issues), strings.Join(e.Issues, "\n  "))
}

// Validate checks an env set and its target overrides for invalid keys or
// empty values that are not explicitly allowed.
func Validate(cfg *config.Config, setName string) error {
	es := cfg.EnvSetByName(setName)
	if es == nil {
		return fmt.Errorf("env set %q not found", setName)
	}

	var issues []string

	for k, v := range es.Vars {
		if !validKeyRe.MatchString(k) {
			issues = append(issues, fmt.Sprintf("base: invalid key name %q", k))
		}
		if strings.TrimSpace(v) == "" {
			issues = append(issues, fmt.Sprintf("base: key %q has empty value", k))
		}
	}

	for _, tgt := range es.Targets {
		for k, v := range tgt.Overrides {
			if !validKeyRe.MatchString(k) {
				issues = append(issues, fmt.Sprintf("target %q: invalid key name %q", tgt.Name, k))
			}
			if strings.TrimSpace(v) == "" {
				issues = append(issues, fmt.Sprintf("target %q: key %q has empty value", tgt.Name, k))
			}
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
