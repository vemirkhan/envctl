package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/envctl/envctl/internal/config"
)

// ResolvedEnv holds the final key-value pairs for a given env set and target.
type ResolvedEnv map[string]string

// Resolve returns the merged environment variables for the given env set name
// and deployment target. Target-specific values override base values.
// If target is empty, only base variables are returned.
func Resolve(cfg *config.Config, setName, target string) (ResolvedEnv, error) {
	set, err := cfg.EnvSetByName(setName)
	if err != nil {
		return nil, fmt.Errorf("resolve: %w", err)
	}

	resolved := make(ResolvedEnv, len(set.Vars))

	for k, v := range set.Vars {
		resolved[k] = expandValue(v)
	}

	if target != "" {
		overlays, ok := set.Targets[target]
		if !ok {
			return nil, fmt.Errorf("resolve: target %q not found in env set %q", target, setName)
		}
		for k, v := range overlays {
			resolved[k] = expandValue(v)
		}
	}

	return resolved, nil
}

// ToExportLines converts a ResolvedEnv into a slice of "export KEY=VALUE" lines
// suitable for sourcing in a shell.
func (r ResolvedEnv) ToExportLines() []string {
	lines := make([]string, 0, len(r))
	for k, v := range r {
		lines = append(lines, fmt.Sprintf("export %s=%s", k, shellQuote(v)))
	}
	return lines
}

// expandValue replaces ${VAR} and $VAR references with values from the current
// process environment.
func expandValue(v string) string {
	return os.ExpandEnv(v)
}

// shellQuote wraps a value in single quotes, escaping any existing single
// quotes so the result is safe to use in POSIX shell.
func shellQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "'\\'')") + "'"
}
