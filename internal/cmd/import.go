package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envctl/internal/config"
	"envctl/internal/env"
)

// NewImportCmd returns the cobra command for importing env vars from a file.
func NewImportCmd() *cobra.Command {
	var (
		format    string
		target    string
		overwrite bool
	)

	cmd := &cobra.Command{
		Use:   "import <set> <file>",
		Short: "Import environment variables from a file into an env set",
		Long: `Import reads variables from a .env or JSON file and merges them
into the base (or a specific target) of the given env set.

Supported formats: dotenv, json`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			file := args[1]

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			n, err := env.Import(cfg, env.ImportOptions{
				File:      file,
				Format:    env.ImportFormat(format),
				SetName:   setName,
				Target:    target,
				Overwrite: overwrite,
			})
			if err != nil {
				return err
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Imported %d variable(s) into %q", n, setName)
			if target != "" {
				fmt.Fprintf(cmd.OutOrStdout(), " (target: %s)", target)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Import format: dotenv or json")
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target to import into (default: base)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys")

	return cmd
}
