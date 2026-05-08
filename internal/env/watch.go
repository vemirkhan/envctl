package env

import (
	"fmt"
	"time"

	"github.com/envctl/envctl/internal/config"
)

// WatchResult holds the outcome of a single watch poll.
type WatchResult struct {
	Set       string
	Target    string
	ChangedAt time.Time
	Diffs     []DiffEntry
}

// WatchOptions configures the Watch poller.
type WatchOptions struct {
	Set      string
	Target   string
	Interval time.Duration
	MaxPolls int // 0 = run forever
}

// Watch polls the resolved env set at the given interval and emits a
// WatchResult whenever the resolved variables change between polls.
// The caller controls termination via the done channel.
func Watch(cfg *config.Config, opts WatchOptions, done <-chan struct{}) (<-chan WatchResult, <-chan error) {
	results := make(chan WatchResult)
	errs := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errs)

		prev, err := Resolve(cfg, opts.Set, opts.Target)
		if err != nil {
			errs <- fmt.Errorf("watch: initial resolve: %w", err)
			return
		}

		polls := 0
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				curr, err := Resolve(cfg, opts.Set, opts.Target)
				if err != nil {
					errs <- fmt.Errorf("watch: poll resolve: %w", err)
					return
				}

				diffs := Diff(prev, curr)
				if len(diffs) > 0 {
					results <- WatchResult{
						Set:       opts.Set,
						Target:    opts.Target,
						ChangedAt: time.Now().UTC(),
						Diffs:     diffs,
					}
					prev = curr
				}

				polls++
				if opts.MaxPolls > 0 && polls >= opts.MaxPolls {
					return
				}
			}
		}
	}()

	return results, errs
}
