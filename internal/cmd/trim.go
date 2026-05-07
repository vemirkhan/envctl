package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/envctl/internal/config"
	"github.com/your-org/envctl/internal/env"
)

// NewTrimCmd returns the cobra command for trimming keys by prefix/suffix.
func NewTrimCmd() *cobra.Command {
	var prefix string
	var suffix string
	var targets []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "trim <set>",
		Short: "Remove keys matching a prefix or suffix from an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			res, err := env.Trim(cfg, setName, prefix, suffix, targets)
			if err != nil {
				return err
			}

			if len(res.Removed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no keys matched — nothing trimmed")
				return nil
			}

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] would remove %d key(s) from %q:\n", len(res.Removed), setName)
				for _, k := range res.Removed {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", k)
				}
				return nil
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "trimmed %d key(s) from %q: %s\n",
				len(res.Removed), setName, strings.Join(res.Removed, ", "))
			return nil
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Remove keys with this prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "Remove keys with this suffix")
	cmd.Flags().StringSliceVar(&targets, "targets", nil, "Also trim matching overrides in these targets")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")

	return cmd
}
