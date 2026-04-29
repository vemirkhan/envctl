package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewSyncCmd returns the cobra command for the `sync` sub-command.
func NewSyncCmd() *cobra.Command {
	var cfgFile string
	var outDir string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync <env-set>",
		Short: "Sync environment variables to target output files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			results, err := env.Sync(cfg, setName, outDir, dryRun)
			if err != nil {
				return err
			}

			for _, r := range results {
				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] would write %d vars to %s (target: %s)\n",
						r.Written, r.File, r.Target)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "synced %d vars → %s (target: %s)\n",
						r.Written, r.File, r.Target)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "c", "envctl.yaml", "path to config file")
	cmd.Flags().StringVarP(&outDir, "out-dir", "o", "", "override output directory for all target files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print actions without writing files")

	_ = os.Getenv // suppress unused import warning in minimal build
	return cmd
}
