package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewSealCmd returns the parent `seal` command with add/remove subcommands.
func NewSealCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seal",
		Short: "Seal or unseal keys in an env set to prevent target overrides",
	}
	cmd.AddCommand(newSealAddCmd())
	cmd.AddCommand(newSealRemoveCmd())
	return cmd
}

func newSealAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <set> [key...]",
		Short: "Seal keys in an env set (omit keys to seal all base keys)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			keys := args[1:]

			res, err := env.Seal(cfg, setName, keys)
			if err != nil {
				return err
			}

			if len(res.Sealed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no new keys sealed")
				return nil
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "sealed in %s: %s\n", setName, strings.Join(res.Sealed, ", "))
			return nil
		},
	}
}

func newSealRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <set> [key...]",
		Short: "Unseal keys in an env set (omit keys to unseal all)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			keys := args[1:]

			res, err := env.Unseal(cfg, setName, keys)
			if err != nil {
				return err
			}

			if len(res.Sealed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no keys unsealed")
				return nil
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "unsealed in %s: %s\n", setName, strings.Join(res.Sealed, ", "))
			return nil
		},
	}
}
