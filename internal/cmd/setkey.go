package cmd

import (
	"fmt"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewSetKeyCmd returns the `setkey` sub-command.
func NewSetKeyCmd() *cobra.Command {
	var (
		cfgPath   string
		target    string
		overwrite bool
	)

	cmd := &cobra.Command{
		Use:   "setkey <env-set> <KEY> <value>",
		Short: "Set or update a single key in an env set",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName, key, value := args[0], args[1], args[2]

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			opts := env.SetKeyOptions{
				Target:    target,
				Overwrite: overwrite,
			}
			if err := env.SetKey(cfg, setName, key, value, opts); err != nil {
				return err
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			loc := "base"
			if target != "" {
				loc = fmt.Sprintf("target %q", target)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "set %s=%s in %s [%s]\n", key, value, setName, loc)
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", "envctl.yaml", "path to config file")
	cmd.Flags().StringVarP(&target, "target", "t", "", "apply to a specific target override instead of base")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing key if present")

	return cmd
}
