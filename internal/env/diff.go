package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// DiffResult holds the comparison between two resolved env sets.
type DiffResult struct {
	Added    map[string]string
	Removed  map[string]string
	Changed  map[string][2]string // key -> [old, new]
	Unchanged map[string]string
}

// Diff compares two resolved environment maps and returns a DiffResult.
func Diff(base, target map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, tv := range target {
		if bv, ok := base[k]; !ok {
			result.Added[k] = tv
		} else if bv != tv {
			result.Changed[k] = [2]string{bv, tv}
		} else {
			result.Unchanged[k] = tv
		}
	}

	for k, bv := range base {
		if _, ok := target[k]; !ok {
			result.Removed[k] = bv
		}
	}

	return result
}

// WriteDiff writes a human-readable diff to w.
func WriteDiff(w io.Writer, d DiffResult) {
	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}

	for _, k := range keys(d.Added) {
		fmt.Fprintf(w, "+ %s=%s\n", k, d.Added[k])
	}
	for _, k := range keys(d.Removed) {
		fmt.Fprintf(w, "- %s=%s\n", k, d.Removed[k])
	}

	changedKeys := make([]string, 0, len(d.Changed))
	for k := range d.Changed {
		changedKeys = append(changedKeys, k)
	}
	sort.Strings(changedKeys)
	for _, k := range changedKeys {
		pair := d.Changed[k]
		fmt.Fprintf(w, "~ %s: %s -> %s\n", k, pair[0], pair[1])
	}

	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		fmt.Fprintln(w, "(no differences)")
	}
	_ = strings.Contains // suppress unused import if needed
}
