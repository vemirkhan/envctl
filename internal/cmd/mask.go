package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewMaskCmd returns the parent mask command with add/remove subcommands.
func NewMaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mask",
		Short: "Mark or unmark environment variable keys as sensitive (masked)",
	}
	cmd.AddCommand(newMaskAddCmd())
	cmd.AddCommand(newMaskRemoveCmd())
	return cmd
}

func newMaskAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <set> [key...]",
		Short: "Mask one or more keys in an env set (omit keys to mask all)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			setName := args[0]
			keys := args[1:]

			res, err := env.Mask(cfg, setName, keys)
			if err != nil {
				return err
			}

			if len(res.Masked) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "masked: %s\n", strings.Join(res.Masked, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "already masked (skipped): %s\n", strings.Join(res.Skipped, ", "))
			}

			return config.Save(cfg, cfgPath)
		},
	}
}

func newMaskRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <set> <key...>",
		Short: "Unmask one or more keys in an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			setName := args[0]
			keys := args[1:]

			res, err := env.Unmask(cfg, setName, keys)
			if err != nil {
				return err
			}

			if len(res.Masked) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "unmasked: %s\n", strings.Join(res.Masked, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "not masked (skipped): %s\n", strings.Join(res.Skipped, ", "))
			}

			return config.Save(cfg, cfgPath)
		},
	}
}
