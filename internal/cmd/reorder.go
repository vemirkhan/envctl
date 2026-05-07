package cmd

import (
	"fmt"
	"strings"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewReorderCmd returns a cobra command that reorders keys within an env set.
func NewReorderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reorder <set> <key1,key2,...>",
		Short: "Reorder keys in an env set",
		Long: `Reorder specifies the desired key order for an env set.
Keys not listed are appended at the end in their current order.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			setName := args[0]
			keys := splitKeys(args[1])
			if len(keys) == 0 {
				return fmt.Errorf("at least one key must be specified")
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if err := env.Reorder(cfg, setName, keys); err != nil {
				return err
			}

			if err := config.Save(cfg, cfgPath); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "reordered keys in %q: %s\n", setName, strings.Join(keys, ", "))
			return nil
		},
	}
	return cmd
}

func splitKeys(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
