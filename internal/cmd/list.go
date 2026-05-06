package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/config"
	"envctl/internal/env"
)

// NewListCmd returns the cobra command for listing env sets.
// The list command displays a table of all env sets defined in the config file,
// showing each set's name, base variable count, and configured targets.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all defined env sets",
		Long:  `Display a table of all env sets with their base variable count and configured targets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return fmt.Errorf("failed to get config flag: %w", err)
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			results := env.List(cfg)
			if len(results) == 0 {
				fmt.Fprintln(os.Stdout, "No env sets defined.")
				return nil
			}

			env.WriteList(os.Stdout, results)
			return nil
		},
	}

	return cmd
}
