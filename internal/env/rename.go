package env

import (
	"fmt"

	"github.com/user/envctl/internal/config"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldName string
	NewName string
	KeysUpdated int
}

// Rename renames an env set from oldName to newName within the config.
// It returns an error if oldName does not exist or newName already exists.
func Rename(cfg *config.Config, oldName, newName string) (*RenameResult, error) {
	if oldName == newName {
		return nil, fmt.Errorf("old name and new name are the same: %q", oldName)
	}

	var oldIdx int = -1
	for i, s := range cfg.EnvSets {
		if s.Name == oldName {
			oldIdx = i
		}
		if s.Name == newName {
			return nil, fmt.Errorf("env set %q already exists", newName)
		}
	}
	if oldIdx == -1 {
		return nil, fmt.Errorf("env set %q not found", oldName)
	}

	keysUpdated := len(cfg.EnvSets[oldIdx].Vars)
	cfg.EnvSets[oldIdx].Name = newName

	// Update any sync targets that reference the old name.
	for i, s := range cfg.EnvSets {
		for j, t := range s.Targets {
			if t.Ref == oldName {
				cfg.EnvSets[i].Targets[j].Ref = newName
			}
		}
	}

	return &RenameResult{
		OldName:     oldName,
		NewName:     newName,
		KeysUpdated: keysUpdated,
	}, nil
}
