package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envctl/internal/config"
	"envctl/internal/env"
)

// NewExportCmd creates the export subcommand which resolves and exports
// environment variables for a given env set and optional deployment target.
func NewExportCmd() *cobra.Command {
	var (
		configPath string
		target     string
		format     string
		output     string
	)

	cmd := &cobra.Command{
		Use:   "export <env-set>",
		Short: "Export environment variables for an env set",
		Long: `Resolve and export environment variables for the specified env set.
Optionally overlay a deployment target on top of the base variables.

Supported formats: export (default), dotenv, json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			vars, err := env.Resolve(cfg, setName, target)
			if err != nil {
				return fmt.Errorf("resolving env set %q: %w", setName, err)
			}

			var w *os.File
			if output == "" || output == "-" {
				w = os.Stdout
			} else {
				w, err = os.Create(output)
				if err != nil {
					return fmt.Errorf("opening output file: %w", err)
				}
				defer w.Close()
			}

			if err := env.Export(vars, format, w); err != nil {
				return fmt.Errorf("exporting vars: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "envctl.yaml", "Path to the envctl config file")
	cmd.Flags().StringVarP(&target, "target", "t", "", "Deployment target overlay (e.g. production, staging)")
	cmd.Flags().StringVarP(&format, "format", "f", "export", "Output format: export, dotenv, json")
	cmd.Flags().StringVarP(&output, "output", "o", "-", "Output file path (default: stdout)")

	return cmd
}
