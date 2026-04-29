package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/config"
	"envctl/internal/env"
)

// NewDiffCmd returns the diff subcommand.
// Usage: envctl diff <env-set> <target-a> <target-b>
func NewDiffCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "diff <env-set> <target-a> <target-b>",
		Short: "Show differences between two deployment targets for an env set",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName, targetA, targetB := args[0], args[1], args[2]

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			resolvedA, err := env.Resolve(cfg, setName, targetA)
			if err != nil {
				return fmt.Errorf("resolving %s/%s: %w", setName, targetA, err)
			}

			resolvedB, err := env.Resolve(cfg, setName, targetB)
			if err != nil {
				return fmt.Errorf("resolving %s/%s: %w", setName, targetB, err)
			}

			fmt.Fprintf(os.Stdout, "Diff %s: %s -> %s\n", setName, targetA, targetB)
			fmt.Fprintln(os.Stdout, "---")

			diffResult := env.Diff(resolvedA, resolvedB)
			env.WriteDiff(os.Stdout, diffResult)

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "envctl.yaml", "path to config file")
	return cmd
}
