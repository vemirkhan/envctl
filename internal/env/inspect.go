package env

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envctl/internal/config"
)

// InspectResult holds the details for a single env set inspection.
type InspectResult struct {
	Name    string
	Base    map[string]string
	Targets map[string]map[string]string
}

// Inspect returns the fully resolved key/value details for the named env set,
// including per-target overrides.
func Inspect(cfg *config.Config, setName string) (*InspectResult, error) {
	set := cfg.EnvSetByName(setName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", setName)
	}

	result := &InspectResult{
		Name:    set.Name,
		Base:    make(map[string]string, len(set.Base)),
		Targets: make(map[string]map[string]string),
	}

	for k, v := range set.Base {
		result.Base[k] = expandValue(v)
	}

	for _, t := range set.Targets {
		overrides := make(map[string]string, len(t.Overrides))
		for k, v := range t.Overrides {
			overrides[k] = expandValue(v)
		}
		result.Targets[t.Name] = overrides
	}

	return result, nil
}

// WriteInspect formats an InspectResult to the provided writer.
func WriteInspect(w io.Writer, r *InspectResult) {
	fmt.Fprintf(w, "Env Set: %s\n", r.Name)
	fmt.Fprintf(w, "\nBase (%d keys):\n", len(r.Base))

	baseKeys := sortedKeys(r.Base)
	for _, k := range baseKeys {
		fmt.Fprintf(w, "  %s=%s\n", k, r.Base[k])
	}

	if len(r.Targets) == 0 {
		return
	}

	targetNames := make([]string, 0, len(r.Targets))
	for name := range r.Targets {
		targetNames = append(targetNames, name)
	}
	sort.Strings(targetNames)

	for _, name := range targetNames {
		overrides := r.Targets[name]
		fmt.Fprintf(w, "\nTarget: %s (%d override(s)):\n", name, len(overrides))
		for _, k := range sortedKeys(overrides) {
			fmt.Fprintf(w, "  %s=%s\n", k, overrides[k])
		}
	}
}
