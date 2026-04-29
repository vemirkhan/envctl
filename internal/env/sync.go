package env

import (
	"fmt"
	"os"

	"github.com/user/envctl/internal/config"
)

// SyncResult holds the outcome of a sync operation for a single target.
type SyncResult struct {
	Target  string
	Written int
	Skipped int
	File    string
}

// Sync writes resolved environment variables for each target defined in the
// env set to its configured output file. If outDir is non-empty it overrides
// the directory portion of each target's path.
func Sync(cfg *config.Config, setName string, outDir string, dryRun bool) ([]SyncResult, error) {
	envSet, err := cfg.EnvSetByName(setName)
	if err != nil {
		return nil, err
	}

	if len(envSet.Targets) == 0 {
		return nil, fmt.Errorf("env set %q has no targets defined", setName)
	}

	var results []SyncResult

	for _, target := range envSet.Targets {
		resolved, err := Resolve(cfg, setName, target.Name)
		if err != nil {
			return results, fmt.Errorf("resolving target %q: %w", target.Name, err)
		}

		outPath := target.File
		if outPath == "" {
			outPath = fmt.Sprintf("%s.env", target.Name)
		}
		if outDir != "" {
			outPath = outDir + "/" + outPath
		}

		format := target.Format
		if format == "" {
			format = "dotenv"
		}

		result := SyncResult{
			Target:  target.Name,
			Written: len(resolved),
			File:    outPath,
		}

		if !dryRun {
			f, err := os.Create(outPath)
			if err != nil {
				return results, fmt.Errorf("creating file %q: %w", outPath, err)
			}
			if err := Export(resolved, format, f); err != nil {
				f.Close()
				return results, fmt.Errorf("exporting to %q: %w", outPath, err)
			}
			f.Close()
		}

		results = append(results, result)
	}

	return results, nil
}
