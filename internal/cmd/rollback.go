package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewRollbackCmd returns the cobra command for rolling back an env set to a snapshot.
func NewRollbackCmd() *cobra.Command {
	var targetName string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "rollback <env-set> <snapshot>",
		Short: "Restore an env set from a named snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			envSetName := args[0]
			snapshotName := args[1]

			res, err := env.Rollback(cfg, envSetName, snapshotName, targetName)
			if err != nil {
				return err
			}

			scope := "base"
			if targetName != "" {
				scope = fmt.Sprintf("target %q", targetName)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Rolled back %s (%s) from snapshot %q — %d key(s) restored:\n",
				res.EnvSet, scope, res.SnapshotName, len(res.Restored))

			keys := make([]string, 0, len(res.Restored))
			for k := range res.Restored {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, res.Restored[k])
			}

			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: changes not saved)")
				return nil
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not save config: %v\n", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&targetName, "target", "t", "", "Restore a specific target's overrides instead of base vars")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview the rollback without saving")
	return cmd
}
