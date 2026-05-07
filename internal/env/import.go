package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envctl/internal/config"
)

// ImportFormat represents the format of the file to import from.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatJSON   ImportFormat = "json"
)

// ImportOptions controls how variables are imported into an env set.
type ImportOptions struct {
	File      string
	Format    ImportFormat
	SetName   string
	Target    string
	Overwrite bool
}

// Import reads environment variables from a file and merges them into the
// specified env set (optionally scoped to a target).
func Import(cfg *config.Config, opts ImportOptions) (int, error) {
	set := cfg.EnvSetByName(opts.SetName)
	if set == nil {
		return 0, fmt.Errorf("env set %q not found", opts.SetName)
	}

	vars, err := parseImportFile(opts.File, opts.Format)
	if err != nil {
		return 0, err
	}

	if opts.Target == "" {
		return applyVars(set.Base, vars, opts.Overwrite)
	}

	targetVars, ok := set.Targets[opts.Target]
	if !ok {
		set.Targets[opts.Target] = make(map[string]string)
		targetVars = set.Targets[opts.Target]
	}
	return applyVars(targetVars, vars, opts.Overwrite)
}

func applyVars(dest, src map[string]string, overwrite bool) (int, error) {
	count := 0
	for k, v := range src {
		if _, exists := dest[k]; exists && !overwrite {
			continue
		}
		dest[k] = v
		count++
	}
	return count, nil
}

func parseImportFile(path string, format ImportFormat) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening import file: %w", err)
	}
	defer f.Close()

	switch format {
	case ImportFormatDotenv:
		return parseDotenv(f)
	case ImportFormatJSON:
		var m map[string]string
		if err := json.NewDecoder(f).Decode(&m); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unsupported import format: %q", format)
	}
}

func parseDotenv(f *os.File) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = val
	}
	return result, scanner.Err()
}
