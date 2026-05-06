package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewTagCmd returns the root `tag` command with add/remove/list subcommands.
func NewTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on env sets",
	}
	cmd.AddCommand(newTagAddCmd())
	cmd.AddCommand(newTagRemoveCmd())
	cmd.AddCommand(newTagListCmd())
	return cmd
}

func newTagAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <set> <tag> [tag...]",
		Short: "Add tags to an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if err := env.Tag(cfg, args[0], args[1:]); err != nil {
				return err
			}
			if err := config.Save(cfg, cfgPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tagged %q with: %s\n", args[0], strings.Join(args[1:], ", "))
			return nil
		},
	}
}

func newTagRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <set> <tag> [tag...]",
		Short: "Remove tags from an env set",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if err := env.Untag(cfg, args[0], args[1:]); err != nil {
				return err
			}
			if err := config.Save(cfg, cfgPath); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed tags from %q: %s\n", args[0], strings.Join(args[1:], ", "))
			return nil
		},
	}
}

func newTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <tag>",
		Short: "List env sets that have a given tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			names := env.ListByTag(cfg, args[0])
			if len(names) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no env sets tagged %q\n", args[0])
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}
