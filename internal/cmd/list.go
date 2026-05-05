package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/config"
	"envctl/internal/env"
)

// NewListCmd returns the cobra command for listing env sets.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all defined env sets",
		Long:  `Display a table of all env sets with their base variable count and configured targets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			results := env.List(cfg)
			env.WriteList(os.Stdout, results)
			return nil
		},
	}

	return cmd
}
