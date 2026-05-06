package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewSnapshotCmd returns the snapshot command with take/list/delete sub-commands.
func NewSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage point-in-time snapshots of env sets",
	}
	cmd.AddCommand(newSnapshotTakeCmd())
	cmd.AddCommand(newSnapshotListCmd())
	cmd.AddCommand(newSnapshotDeleteCmd())
	return cmd
}

func newSnapshotTakeCmd() *cobra.Command {
	var target, name string
	cmd := &cobra.Command{
		Use:   "take <set>",
		Short: "Capture a snapshot of an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			snap, err := env.TakeSnapshot(cfg, args[0], target, name)
			if err != nil {
				return err
			}
			if err := config.Save(cfgPath, cfg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %q created (%d vars)\n", snap.Name, len(snap.Vars))
			return nil
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", "", "target override")
	cmd.Flags().StringVarP(&name, "name", "n", "", "snapshot name (auto-generated if omitted)")
	return cmd
}

func newSnapshotListCmd() *cobra.Command {
	var setFilter string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			snaps := env.ListSnapshots(cfg, setFilter)
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSET\tTARGET\tCREATED")
			for _, s := range snaps {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.Set, s.Target, s.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}
	cmd.Flags().StringVarP(&setFilter, "set", "s", "", "filter by env set name")
	return cmd
}

func newSnapshotDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a snapshot by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if err := env.DeleteSnapshot(cfg, args[0]); err != nil {
				return err
			}
			if err := config.Save(cfgPath, cfg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %q deleted\n", args[0])
			return nil
		},
	}
}

// ensure os import used via potential future file writes
var _ = os.Stderr
