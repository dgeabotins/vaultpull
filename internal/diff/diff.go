package diff

import "fmt"

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string
	Status string // "added", "removed", "changed", "unchanged"
	OldVal string
	NewVal string
}

// Compare returns a diff between old and new secret maps.
func Compare(old, new map[string]string) []Result {
	seen := map[string]bool{}
	var results []Result

	for k, newVal := range new {
		seen[k] = true
		if oldVal, ok := old[k]; !ok {
			results = append(results, Result{Key: k, Status: "added", NewVal: newVal})
		} else if oldVal != newVal {
			results = append(results, Result{Key: k, Status: "changed", OldVal: oldVal, NewVal: newVal})
		} else {
			results = append(results, Result{Key: k, Status: "unchanged", OldVal: oldVal, NewVal: newVal})
		}
	}

	for k, oldVal := range old {
		if !seen[k] {
			results = append(results, Result{Key: k, Status: "removed", OldVal: oldVal})
		}
	}

	return results
}

// Summary returns a human-readable summary string.
func Summary(results []Result) string {
	added, changed, removed := 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case "added":
			added++
		case "changed":
			changed++
		case "removed":
			removed++
		}
	}
	return fmt.Sprintf("+%d added, ~%d changed, -%d removed", added, changed, removed)
}
