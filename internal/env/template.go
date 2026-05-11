package env

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/envctl/envctl/internal/config"
)

// TemplateResult holds the rendered output of a template.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

var placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// Template renders a text template by substituting {{KEY}} placeholders
// with values from the resolved env set. Missing keys are collected and
// returned rather than causing a hard error, unless strict is true.
func Template(cfg *config.Config, setName, target, tmpl string, strict bool) (*TemplateResult, error) {
	vars, err := Resolve(cfg, setName, target)
	if err != nil {
		return nil, err
	}

	missing := []string{}
	seen := map[string]bool{}

	rendered := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := strings.TrimSpace(placeholderRe.FindStringSubmatch(match)[1])
		val, ok := vars[key]
		if !ok {
			if !seen[key] {
				missing = append(missing, key)
				seen[key] = true
			}
			return match
		}
		return val
	})

	sort.Strings(missing)

	if strict && len(missing) > 0 {
		return nil, fmt.Errorf("template: unresolved placeholders: %s", strings.Join(missing, ", "))
	}

	return &TemplateResult{Rendered: rendered, Missing: missing}, nil
}
