package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewValidateCmd returns a cobra command that validates an env set.
// It checks that all keys in the named set are present and that their
// values conform to any constraints defined in the config file.
func NewValidateCmd() *cobra.Command {
	var cfgFile string

	cmd := &cobra.Command{
		Use:   "validate <env-set>",
		Short: "Validate keys and values in an env set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setName := args[0]
			if err := env.Validate(cfg, setName); err != nil {
				fmt.Fprintf(os.Stderr, "envctl: %v\n", err)
				return fmt.Errorf("validating env set %q: %w", setName, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "env set %q is valid\n", setName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "c", "envctl.yaml", "path to config file")
	return cmd
}
