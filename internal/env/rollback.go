package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// RollbackResult describes the outcome of a rollback operation.
type RollbackResult struct {
	EnvSet       string
	SnapshotName string
	Restored     map[string]string
}

// Rollback restores an env set's base vars from a named snapshot.
// If targetName is non-empty, only that target's overrides are restored.
func Rollback(cfg *config.Config, envSetName, snapshotName, targetName string) (*RollbackResult, error) {
	set := cfg.EnvSetByName(envSetName)
	if set == nil {
		return nil, fmt.Errorf("env set %q not found", envSetName)
	}

	var snap *config.Snapshot
	for i := range cfg.Snapshots {
		if cfg.Snapshots[i].EnvSet == envSetName && cfg.Snapshots[i].Name == snapshotName {
			snap = &cfg.Snapshots[i]
			break
		}
	}
	if snap == nil {
		return nil, fmt.Errorf("snapshot %q not found for env set %q", snapshotName, envSetName)
	}

	restored := make(map[string]string)

	if targetName == "" {
		// Restore base vars
		for k, v := range snap.Vars {
			set.Base[k] = v
			restored[k] = v
		}
	} else {
		// Restore a specific target's overrides
		snapTarget, ok := snap.Targets[targetName]
		if !ok {
			return nil, fmt.Errorf("target %q not found in snapshot %q", targetName, snapshotName)
		}
		for i := range set.Targets {
			if set.Targets[i].Name == targetName {
				set.Targets[i].Overrides = make(map[string]string)
				for k, v := range snapTarget {
					set.Targets[i].Overrides[k] = v
					restored[k] = v
				}
				break
			}
		}
	}

	return &RollbackResult{
		EnvSet:       envSetName,
		SnapshotName: snapshotName,
		Restored:     restored,
	}, nil
}
