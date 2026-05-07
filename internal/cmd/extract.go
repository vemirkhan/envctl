package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewExtractCmd returns the cobra command for extracting keys into a new env set.
func NewExtractCmd() *cobra.Command {
	var target string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "extract <set> <dest> <KEY1,KEY2,...>",
		Short: "Extract specific keys from an env set into a new set",
		Long: `Extract copies a subset of resolved keys from <set> (optionally merged
with <target>) into a brand-new env set named <dest>.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			srcName := args[0]
			destName := args[1]
			keys := splitExtractKeys(args[2])

			if len(keys) == 0 {
				return fmt.Errorf("at least one key must be specified")
			}

			res, err := env.Extract(cfg, srcName, keys, destName, target, overwrite)
			if err != nil {
				return err
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Extracted %d key(s) from %q into %q\n",
				len(res.Keys), srcName, res.SetName)
			for _, k := range res.Keys {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, res.Vars[k])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Resolve keys using this target override")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination set if it already exists")
	return cmd
}

func splitExtractKeys(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if k := strings.TrimSpace(p); k != "" {
			out = append(out, k)
		}
	}
	return out
}
