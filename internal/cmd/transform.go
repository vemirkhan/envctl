package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewTransformCmd returns the cobra command for the transform subcommand.
func NewTransformCmd() *cobra.Command {
	var target string
	var keys string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "transform <set> <op>",
		Short: "Apply a value transformation to keys in an env set",
		Long: `Apply a built-in transformation to key values within an env set.

Supported ops: upper, lower, trim, base64

Examples:
  envctl transform app upper --keys APP_NAME,REGION
  envctl transform app trim --target prod --dry-run`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if t := strings.TrimSpace(k); t != "" {
						keyList = append(keyList, t)
					}
				}
			}

			changed, err := env.Transform(cfg, env.TransformOptions{
				SetName: args[0],
				Op:      env.TransformOp(args[1]),
				Target:  target,
				Keys:    keyList,
				DryRun:  dryRun,
			})
			if err != nil {
				return err
			}

			if len(changed) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no values changed")
				return nil
			}

			for k, v := range changed {
				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] %s -> %s\n", k, v)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "transformed %s -> %s\n", k, v)
				}
			}

			if !dryRun {
				if err := config.Save(cfg, cfgPath); err != nil {
					return fmt.Errorf("save config: %w", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "apply transform to a specific target instead of base")
	cmd.Flags().StringVarP(&keys, "keys", "k", "", "comma-separated list of keys to transform (default: all)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")

	_ = os.Stderr // suppress unused import
	return cmd
}
