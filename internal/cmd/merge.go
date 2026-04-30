package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewMergeCmd returns the cobra command for merging env sets.
func NewMergeCmd() *cobra.Command {
	var cfgPath string
	var format string
	var showConflicts bool

	cmd := &cobra.Command{
		Use:   "merge <set1> <set2> [set3...]",
		Short: "Merge multiple env sets into a single variable map",
		Long: `Merge combines two or more env sets into one.
Later sets take precedence over earlier ones.
Use --conflicts to print a summary of overridden keys.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgPath == "" {
				cfgPath, _ = cmd.Flags().GetString("config")
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			result, err := env.Merge(cfg, args)
			if err != nil {
				return err
			}

			if showConflicts && len(result.Conflicts) > 0 {
				fmt.Fprintln(os.Stderr, "# Conflicts (last set wins):")
				for _, c := range result.Conflicts {
					parts := make([]string, 0, len(c.Values))
					for set, val := range c.Values {
						parts = append(parts, fmt.Sprintf("%s=%q", set, val))
					}
					fmt.Fprintf(os.Stderr, "#   %s: %s\n", c.Key, strings.Join(parts, ", "))
				}
			}

			return env.Export(result.Vars, format, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", "envctl.yaml", "path to config file")
	cmd.Flags().StringVarP(&format, "format", "f", "export", "output format: export, dotenv, json")
	cmd.Flags().BoolVar(&showConflicts, "conflicts", false, "print conflicting keys to stderr")

	return cmd
}
