package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// CopyResult holds the result of a copy operation.
type CopyResult struct {
	SourceSet string
	DestSet   string
	Keys      []string
}

// Copy duplicates the resolved variables from one env set to another within
// the config, merging overrides from an optional target. It returns the list
// of keys that were copied.
func Copy(cfg *config.Config, srcName, dstName, target string) (*CopyResult, error) {
	src, err := cfg.EnvSetByName(srcName)
	if err != nil {
		return nil, fmt.Errorf("copy: source: %w", err)
	}

	dst, err := cfg.EnvSetByName(dstName)
	if err != nil {
		return nil, fmt.Errorf("copy: destination: %w", err)
	}

	srcResolved, err := Resolve(cfg, srcName, target)
	if err != nil {
		return nil, fmt.Errorf("copy: resolve source: %w", err)
	}

	if dst.Base == nil {
		dst.Base = make(map[string]string)
	}

	copied := make([]string, 0, len(srcResolved))
	for k, v := range srcResolved {
		dst.Base[k] = v
		copied = append(copied, k)
	}

	_ = src // src used for validation above

	return &CopyResult{
		SourceSet: srcName,
		DestSet:   dstName,
		Keys:      sortedKeys(srcResolved),
	}, nil
}
