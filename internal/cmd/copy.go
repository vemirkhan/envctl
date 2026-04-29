package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envctl/internal/config"
	"github.com/user/envctl/internal/env"
)

// NewCopyCmd returns the cobra command for copying env sets.
func NewCopyCmd() *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "copy <src-env-set> <dst-env-set>",
		Short: "Copy resolved variables from one env set into another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			src, dst := args[0], args[1]
			result, err := env.Copy(cfg, src, dst, target)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(),
				"Copied %d key(s) from %q to %q: %s\n",
				len(result.Keys),
				result.SourceSet,
				result.DestSet,
				strings.Join(result.Keys, ", "),
			)
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Optional deployment target to include overrides from")
	return cmd
}
