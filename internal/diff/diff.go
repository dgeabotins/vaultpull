package diff

import "fmt"

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Changed ChangeType = "changed"
	Removed ChangeType = "removed"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level difference.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare returns the list of changes between old and new secret maps.
func Compare(oldSecrets, newSecrets map[string]string) []Change {
	var changes []Change

	for k, newVal := range newSecrets {
		oldVal, exists := oldSecrets[k]
		if !exists {
			changes = append(changes, Change{Key: k, Type: Added, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: Changed, OldValue: oldVal, NewValue: newVal})
		} else {
			changes = append(changes, Change{Key: k, Type: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range oldSecrets {
		if _, exists := newSecrets[k]; !exists {
			changes = append(changes, Change{Key: k, Type: Removed, OldValue: oldVal})
		}
	}

	return changes
}

// Summary returns a human-readable summary of changes.
func Summary(changes []Change) string {
	added, changed, removed := 0, 0, 0
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Changed:
			changed++
		case Removed:
			removed++
		}
	}
	return fmt.Sprintf("+%d added, ~%d changed, -%d removed", added, changed, removed)
}
