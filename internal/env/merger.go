package env

import (
	"os"
	"strings"
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added   int
	Updated int
	Unchanged int
}

// MergeIntoFile merges newSecrets into an existing env file at path.
// Existing keys are updated; new keys are appended. Comments are preserved.
func MergeIntoFile(path string, newSecrets map[string]string) (MergeResult, error) {
	existing, err := LoadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return MergeResult{}, err
	}

	result := MergeResult{}
	merged := make(map[string]string)

	// Copy existing
	for k, v := range existing {
		merged[k] = v
	}

	for k, v := range newSecrets {
		old, exists := existing[k]
		if !exists {
			result.Added++
		} else if old != v {
			result.Updated++
		} else {
			result.Unchanged++
		}
		merged[k] = v
	}

	if err := WriteFile(path, merged); err != nil {
		return MergeResult{}, err
	}
	return result, nil
}

// Summary returns a human-readable summary of the merge result.
func (r MergeResult) Summary() string {
	parts := []string{}
	if r.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", r.Added))
	}
	if r.Updated > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", r.Updated))
	}
	if r.Unchanged > 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", r.Unchanged))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}
