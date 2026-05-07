package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root cobra command with all subcommands registered.
func NewRootCmd() *cobra.Command {
	var cfgPath string

	root := &cobra.Command{
		Use:   "envctl",
		Short: "Manage and sync environment variable sets across projects",
		Long: `envctl lets you define, resolve, export, and sync environment
variable sets across multiple projects and deployment targets
from a single YAML configuration file.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVarP(&cfgPath, "config", "c", "envctl.yaml", "Path to config file")

	root.AddCommand(
		NewExportCmd(),
		NewSyncCmd(),
		NewValidateCmd(),
		NewDiffCmd(),
		NewCopyCmd(),
		NewMergeCmd(),
		NewRenameCmd(),
		NewListCmd(),
		NewInspectCmd(),
		NewSnapshotCmd(),
		NewCompareCmd(),
		NewTagCmd(),
		NewAuditCmd(),
		NewPromoteCmd(),
		NewRollbackCmd(),
		NewPinCmd(),
		NewSealCmd(),
		NewReorderCmd(),
		NewImportCmd(),
	)

	return root
}

// Execute runs the root command.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
