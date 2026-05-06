package env

import (
	"fmt"
	"io"
	"sort"

	"github.com/envctl/envctl/internal/config"
)

// CompareResult holds the result of comparing two env sets across targets.
type CompareResult struct {
	SetA    string
	SetB    string
	OnlyInA map[string]string
	OnlyInB map[string]string
	Same    map[string]string
	Differ  map[string][2]string
}

// Compare compares two env sets (optionally at a given target) and returns a CompareResult.
func Compare(cfg *config.Config, setA, setB, target string) (*CompareResult, error) {
	resolvedA, err := Resolve(cfg, setA, target)
	if err != nil {
		return nil, fmt.Errorf("resolving %q: %w", setA, err)
	}
	resolvedB, err := Resolve(cfg, setB, target)
	if err != nil {
		return nil, fmt.Errorf("resolving %q: %w", setB, err)
	}

	result := &CompareResult{
		SetA:   setA,
		SetB:   setB,
		OnlyInA: make(map[string]string),
		OnlyInB: make(map[string]string),
		Same:   make(map[string]string),
		Differ: make(map[string][2]string),
	}

	for k, vA := range resolvedA {
		if vB, ok := resolvedB[k]; ok {
			if vA == vB {
				result.Same[k] = vA
			} else {
				result.Differ[k] = [2]string{vA, vB}
			}
		} else {
			result.OnlyInA[k] = vA
		}
	}
	for k, vB := range resolvedB {
		if _, ok := resolvedA[k]; !ok {
			result.OnlyInB[k] = vB
		}
	}
	return result, nil
}

// WriteCompare writes a human-readable comparison to w.
func WriteCompare(w io.Writer, r *CompareResult) {
	fmt.Fprintf(w, "Comparing %q vs %q\n", r.SetA, r.SetB)

	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(w, "\nOnly in %q:\n", r.SetA)
		for _, k := range keys(r.OnlyInA) {
			fmt.Fprintf(w, "  + %s=%s\n", k, r.OnlyInA[k])
		}
	}
	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(w, "\nOnly in %q:\n", r.SetB)
		for _, k := range keys(r.OnlyInB) {
			fmt.Fprintf(w, "  + %s=%s\n", k, r.OnlyInB[k])
		}
	}
	if len(r.Differ) > 0 {
		dkeys := make([]string, 0, len(r.Differ))
		for k := range r.Differ {
			dkeys = append(dkeys, k)
		}
		sort.Strings(dkeys)
		fmt.Fprintln(w, "\nDifferences:")
		for _, k := range dkeys {
			pair := r.Differ[k]
			fmt.Fprintf(w, "  ~ %s: %q → %q\n", k, pair[0], pair[1])
		}
	}
	if len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.Differ) == 0 {
		fmt.Fprintln(w, "No differences found.")
	}
}
