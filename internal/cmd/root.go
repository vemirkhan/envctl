package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// NewRootCmd builds and returns the root cobra command for envctl.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envctl",
		Short: "Manage and sync environment variable sets across projects",
		Long: `envctl lets you define, validate, export, diff, and sync
environment variable sets across multiple projects and deployment
targets from a single YAML config file.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c",
		"envctl.yaml",
		"path to the envctl config file",
	)

	root.AddCommand(NewExportCmd())
	root.AddCommand(NewDiffCmd())
	root.AddCommand(NewSyncCmd())
	root.AddCommand(NewValidateCmd())

	return root
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		_, _ = os.Stderr.WriteString("error: " + err.Error() + "\n")
		os.Exit(1)
	}
}
