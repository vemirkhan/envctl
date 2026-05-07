package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewProtectCmd returns the root protect command with subcommands.
func NewProtectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "protect",
		Short: "Manage protected keys in an env set",
	}
	cmd.AddCommand(newProtectAddCmd())
	cmd.AddCommand(newProtectRemoveCmd())
	cmd.AddCommand(newProtectListCmd())
	return cmd
}

func newProtectAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <set> <key> [key...]",
		Short: "Mark keys as protected in an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if err := env.Protect(cfg, args[0], args[1:]); err != nil {
				return err
			}
			return config.Save(cfg, cfgPath)
		},
	}
}

func newProtectRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <set> <key> [key...]",
		Short: "Remove protection from keys in an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if err := env.Unprotect(cfg, args[0], args[1:]); err != nil {
				return err
			}
			return config.Save(cfg, cfgPath)
		},
	}
}

func newProtectListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <set>",
		Short: "List protected keys in an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			keys, err := env.ProtectedKeys(cfg, args[0])
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no protected keys)")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(keys, "\n"))
			return nil
		},
	}
}
