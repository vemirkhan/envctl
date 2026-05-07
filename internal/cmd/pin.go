package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewPinCmd returns the root `pin` command with `add` and `remove` subcommands.
func NewPinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin or unpin keys in an env set to protect them from sync/promote",
	}
	cmd.AddCommand(newPinAddCmd())
	cmd.AddCommand(newPinRemoveCmd())
	return cmd
}

func newPinAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <set> <key> [key...]",
		Short: "Pin one or more keys in an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			setName, keys := args[0], args[1:]
			results, err := env.Pin(cfg, setName, keys)
			if err != nil {
				return err
			}
			for _, r := range results {
				if r.Skipped {
					fmt.Fprintf(cmd.OutOrStdout(), "skip  %s  (%s)\n", r.Key, r.Reason)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "pinned  %s\n", r.Key)
				}
			}
			pinned := pinnedKeys(cfg, setName)
			fmt.Fprintf(cmd.OutOrStdout(), "pinned keys in %q: [%s]\n", setName, strings.Join(pinned, ", "))
			return config.Save(cfg, cfgPath)
		},
	}
}

func newPinRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <set> <key> [key...]",
		Short: "Unpin one or more keys in an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			setName, keys := args[0], args[1:]
			results, err := env.Unpin(cfg, setName, keys)
			if err != nil {
				return err
			}
			for _, r := range results {
				if r.Skipped {
					fmt.Fprintf(cmd.OutOrStdout(), "skip  %s  (%s)\n", r.Key, r.Reason)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "unpinned  %s\n", r.Key)
				}
			}
			return config.Save(cfg, cfgPath)
		},
	}
}

func pinnedKeys(cfg *config.Config, setName string) []string {
	set := cfg.EnvSetByName(setName)
	if set == nil || len(set.Pinned) == 0 {
		return []string{}
	}
	return set.Pinned
}
