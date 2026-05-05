package env

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"envctl/internal/config"
)

// ListResult holds summary info for a single env set.
type ListResult struct {
	Name    string
	BaseLen int
	Targets []string
}

// List returns a summary of all env sets defined in the config.
func List(cfg *config.Config) []ListResult {
	results := make([]ListResult, 0, len(cfg.EnvSets))

	for _, es := range cfg.EnvSets {
		targetNames := make([]string, 0, len(es.Targets))
		for _, t := range es.Targets {
			targetNames = append(targetNames, t.Name)
		}
		sort.Strings(targetNames)

		results = append(results, ListResult{
			Name:    es.Name,
			BaseLen: len(es.Base),
			Targets: targetNames,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}

// WriteList writes a human-readable table of env sets to w.
func WriteList(w io.Writer, results []ListResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No env sets defined.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tBASE VARS\tTARGETS")
	fmt.Fprintln(tw, "----\t---------\t-------")

	for _, r := range results {
		targets := "-"
		if len(r.Targets) > 0 {
			targets = ""
			for i, t := range r.Targets {
				if i > 0 {
					targets += ", "
				}
				targets += t
			}
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\n", r.Name, r.BaseLen, targets)
	}

	tw.Flush()
}
