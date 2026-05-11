package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewTemplateCmd returns the `template` subcommand.
func NewTemplateCmd() *cobra.Command {
	var target string
	var strict bool
	var file string

	cmd := &cobra.Command{
		Use:   "template <env-set> <template-string>",
		Short: "Render a text template using an env set's resolved variables",
		Long: `Render a template string (or file) by substituting {{KEY}} placeholders
with values from the specified env set, optionally merged with a target.`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			setName := args[0]

			var tmpl string
			if file != "" {
				bytes, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("read template file: %w", err)
				}
				tmpl = string(bytes)
			} else if len(args) == 2 {
				tmpl = args[1]
			} else {
				return fmt.Errorf("provide a template string or --file")
			}

			res, err := env.Template(cfg, setName, target, tmpl, strict)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), res.Rendered)

			if len(res.Missing) > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: unresolved placeholders: %s\n",
					strings.Join(res.Missing, ", "))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Target overlay to apply")
	cmd.Flags().BoolVar(&strict, "strict", false, "Fail if any placeholder is unresolved")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a template file (overrides inline template arg)")

	return cmd
}
