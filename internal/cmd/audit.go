package cmd

import (
	"fmt"
	"os"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewAuditCmd returns the cobra command for the audit subcommand.
func NewAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit <env-set>",
		Short: "Audit an env set for common issues",
		Long: `Audit inspects an env set and reports:
  - Redundant target overrides (same value as base)
  - Target keys not present in base
  - Empty values in base or target overrides`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			if cfgPath == "" {
				cfgPath, _ = cmd.Root().PersistentFlags().GetString("config")
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			results, err := env.Audit(cfg, setName)
			if err != nil {
				return err
			}

			env.WriteAudit(os.Stdout, results)

			// Exit with non-zero if any issues were found.
			for _, r := range results {
				if len(r.UnusedOverrides) > 0 || len(r.MissingInBase) > 0 || len(r.EmptyValues) > 0 {
					os.Exit(1)
				}
			}
			return nil
		},
	}
	return cmd
}
