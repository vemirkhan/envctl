package cmd

import (
	"fmt"
	"os"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewCompareCmd returns a cobra command that compares two env sets.
func NewCompareCmd() *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "compare <set-a> <set-b>",
		Short: "Compare two env sets side by side",
		Long:  `Compare two named env sets, optionally at a specific deployment target.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			if cfgPath == "" {
				cfgPath, _ = cmd.Root().PersistentFlags().GetString("config")
			}
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			setA := args[0]
			setB := args[1]

			result, err := env.Compare(cfg, setA, setB, target)
			if err != nil {
				return err
			}

			env.WriteCompare(os.Stdout, result)

			if len(result.Differ) > 0 || len(result.OnlyInA) > 0 || len(result.OnlyInB) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Deployment target to resolve overrides for")
	return cmd
}
