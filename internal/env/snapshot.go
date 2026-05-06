package env

import (
	"fmt"
	"sort"
	"time"

	"github.com/envctl/envctl/internal/config"
)

// Snapshot captures the resolved state of an env set at a point in time.
type Snapshot struct {
	Name      string            `json:"name"`
	Set       string            `json:"set"`
	Target    string            `json:"target,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	Vars      map[string]string `json:"vars"`
}

// TakeSnapshot resolves the given env set (and optional target) and records it
// as a named snapshot inside the config.
func TakeSnapshot(cfg *config.Config, setName, target, snapshotName string) (*Snapshot, error) {
	vars, err := Resolve(cfg, setName, target)
	if err != nil {
		return nil, fmt.Errorf("snapshot: resolve failed: %w", err)
	}

	if snapshotName == "" {
		snapshotName = fmt.Sprintf("%s-%d", setName, time.Now().UnixMilli())
	}

	for _, s := range cfg.Snapshots {
		if s.Name == snapshotName {
			return nil, fmt.Errorf("snapshot %q already exists", snapshotName)
		}
	}

	snap := config.SnapshotEntry{
		Name:      snapshotName,
		Set:       setName,
		Target:    target,
		CreatedAt: time.Now().UTC(),
		Vars:      vars,
	}
	cfg.Snapshots = append(cfg.Snapshots, snap)

	return &Snapshot{
		Name:      snap.Name,
		Set:       snap.Set,
		Target:    snap.Target,
		CreatedAt: snap.CreatedAt,
		Vars:      snap.Vars,
	}, nil
}

// ListSnapshots returns all snapshots, optionally filtered by set name.
func ListSnapshots(cfg *config.Config, setFilter string) []config.SnapshotEntry {
	result := make([]config.SnapshotEntry, 0, len(cfg.Snapshots))
	for _, s := range cfg.Snapshots {
		if setFilter == "" || s.Set == setFilter {
			result = append(result, s)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

// DeleteSnapshot removes a snapshot by name.
func DeleteSnapshot(cfg *config.Config, name string) error {
	for i, s := range cfg.Snapshots {
		if s.Name == name {
			cfg.Snapshots = append(cfg.Snapshots[:i], cfg.Snapshots[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("snapshot %q not found", name)
}
