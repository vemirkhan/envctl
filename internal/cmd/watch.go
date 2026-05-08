package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envctl/envctl/internal/config"
	"github.com/envctl/envctl/internal/env"
	"github.com/spf13/cobra"
)

// NewWatchCmd returns the watch subcommand.
func NewWatchCmd() *cobra.Command {
	var target string
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "watch <env-set>",
		Short: "Poll an env set and print diffs when values change",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			setName := args[0]
			opts := env.WatchOptions{
				Set:      setName,
				Target:   target,
				Interval: interval,
			}

			done := make(chan struct{})
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sig
				close(done)
			}()

			fmt.Fprintf(cmd.OutOrStdout(), "watching %q (interval %s) — press Ctrl+C to stop\n", setName, interval)

			results, errs := env.Watch(cfg, opts, done)
			for r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] change detected in %q:\n", r.ChangedAt.Format(time.RFC3339), r.Set)
				env.WriteDiff(cmd.OutOrStdout(), r.Diffs)
			}
			if err := <-errs; err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "deployment target overlay")
	cmd.Flags().DurationVarP(&interval, "interval", "i", 5*time.Second, "poll interval")
	return cmd
}
