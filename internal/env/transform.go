package env

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
)

// TransformOp defines a transformation operation to apply to env values.
type TransformOp string

const (
	TransformUpper  TransformOp = "upper"
	TransformLower  TransformOp = "lower"
	TransformTrimWS TransformOp = "trim"
	TransformBase64 TransformOp = "base64"
)

// TransformOptions controls how Transform behaves.
type TransformOptions struct {
	SetName string
	Target  string
	Keys    []string
	Op      TransformOp
	DryRun  bool
}

// Transform applies a named operation to the values of specified keys (or all
// keys if none are specified) within an env set.
func Transform(cfg *config.Config, opts TransformOptions) (map[string]string, error) {
	set := cfg.EnvSetByName(opts.SetName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", opts.SetName)
	}

	if !validOp(opts.Op) {
		return nil, fmt.Errorf("unknown transform op %q: choose upper, lower, trim, or base64", opts.Op)
	}

	var vars map[string]string
	if opts.Target == "" {
		vars = set.Base
	} else {
		tgt, ok := set.Targets[opts.Target]
		if !ok {
			return nil, fmt.Errorf("target %q not found in set %q", opts.Target, opts.SetName)
		}
		vars = tgt
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range vars {
			keys = append(keys, k)
		}
	}

	changed := make(map[string]string)
	for _, k := range keys {
		v, ok := vars[k]
		if !ok {
			return nil, fmt.Errorf("key %q not found", k)
		}
		nv := applyOp(opts.Op, v)
		if nv != v {
			changed[k] = nv
			if !opts.DryRun {
				vars[k] = nv
			}
		}
	}
	return changed, nil
}

func validOp(op TransformOp) bool {
	switch op {
	case TransformUpper, TransformLower, TransformTrimWS, TransformBase64:
		return true
	}
	return false
}

func applyOp(op TransformOp, v string) string {
	switch op {
	case TransformUpper:
		return strings.ToUpper(v)
	case TransformLower:
		return strings.ToLower(v)
	case TransformTrimWS:
		return strings.TrimSpace(v)
	case TransformBase64:
		return encodeBase64(v)
	}
	return v
}
