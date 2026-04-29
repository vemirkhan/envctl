package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	FormatExport Format = "export"
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// Export writes the resolved environment variables to w in the given format.
func Export(w io.Writer, vars map[string]string, format Format) error {
	keys := sortedKeys(vars)

	switch format {
	case FormatExport:
		return writeExport(w, keys, vars)
	case FormatDotenv:
		return writeDotenv(w, keys, vars)
	case FormatJSON:
		return writeJSON(w, keys, vars)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func writeExport(w io.Writer, keys []string, vars map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, shellQuote(vars[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeDotenv(w io.Writer, keys []string, vars map[string]string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, shellQuote(vars[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, keys []string, vars map[string]string) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("  %q: %q%s\n", k, vars[k], comma))
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
