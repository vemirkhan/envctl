package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewPromoteCmd returns the cobra command for promoting target overrides.
func NewPromoteCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "promote <env-set> <from-target> <to-target>",
		Short: "Promote target overrides from one target to another",
		Long: `Copy target-level overrides from one deployment target to another
within the same env set. Existing keys in the destination are preserved
unless --overwrite is specified.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			result, err := env.Promote(cfg, args[0], args[1], args[2], overwrite)
			if err != nil {
				return err
			}

			if len(result.KeysPromoted) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No keys promoted from %q to %q (all keys already present).\n",
					result.FromTarget, result.ToTarget)
				return nil
			}

			sort.Strings(result.KeysPromoted)
			fmt.Fprintf(cmd.OutOrStdout(), "Promoted %d key(s) from %q to %q in env set %q:\n",
				len(result.KeysPromoted), result.FromTarget, result.ToTarget, result.EnvSet)
			for _, k := range result.KeysPromoted {
				fmt.Fprintf(cmd.OutOrStdout(), "  + %s\n", k)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in the destination target")
	return cmd
}
