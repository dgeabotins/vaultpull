package env

import (
	"fmt"
	"os"
	"strings"
)

// PatchOp represents a single key operation.
type PatchOp struct {
	Key    string
	Value  string
	Delete bool
}

// PatchResult summarizes the outcome of a patch operation.
type PatchResult struct {
	Set     []string
	Deleted []string
}

func (r PatchResult) Summary() string {
	var parts []string
	if len(r.Set) > 0 {
		parts = append(parts, fmt.Sprintf("%d set", len(r.Set)))
	}
	if len(r.Deleted) > 0 {
		parts = append(parts, fmt.Sprintf("%d deleted", len(r.Deleted)))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// Patch applies a list of PatchOps to an env file, creating it if needed.
func Patch(path string, ops []PatchOp) (PatchResult, error) {
	data, err := LoadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return PatchResult{}, fmt.Errorf("patch: read %s: %w", path, err)
	}
	if data == nil {
		data = []EnvEntry{}
	}

	m := make(map[string]string)
	for _, e := range data {
		m[e.Key] = e.Value
	}

	var result PatchResult
	for _, op := range ops {
		if op.Delete {
			if _, ok := m[op.Key]; ok {
				delete(m, op.Key)
				result.Deleted = append(result.Deleted, op.Key)
			}
		} else {
			m[op.Key] = op.Value
			result.Set = append(result.Set, op.Key)
		}
	}

	if err := WriteFile(path, m, 0600); err != nil {
		return PatchResult{}, fmt.Errorf("patch: write %s: %w", path, err)
	}
	return result, nil
}
