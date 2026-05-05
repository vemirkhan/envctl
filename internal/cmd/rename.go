package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewRenameCmd returns the cobra command for renaming an env set.
func NewRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old-name> <new-name>",
		Short: "Rename an environment variable set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			if cfgPath == "" {
				cfgPath, _ = cmd.Root().PersistentFlags().GetString("config")
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			oldName := args[0]
			newName := args[1]

			result, err := env.Rename(cfg, oldName, newName)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "renamed env set %q -> %q (%d vars)\n",
				result.OldName, result.NewName, result.KeysUpdated)
			return nil
		},
	}
	return cmd
}
