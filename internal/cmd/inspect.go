package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewInspectCmd returns the cobra command for inspecting an env set.
func NewInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect <env-set>",
		Short: "Show all keys and target overrides for an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			if cfgPath == "" {
				cfgPath, err = cmd.Root().PersistentFlags().GetString("config")
				if err != nil {
					return err
				}
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			setName := args[0]
			result, err := env.Inspect(cfg, setName)
			if err != nil {
				return err
			}

			env.WriteInspect(os.Stdout, result)
			return nil
		},
	}

	return cmd
}
