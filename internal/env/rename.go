package env

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
	Found  bool
}

func (r RenameResult) Summary() string {
	if !r.Found {
		return fmt.Sprintf("key %q not found", r.OldKey)
	}
	return fmt.Sprintf("renamed %q -> %q", r.OldKey, r.NewKey)
}

// Rename loads the env file, renames oldKey to newKey, and writes the result back.
// Returns an error if the file cannot be read or written.
func Rename(path, oldKey, newKey string) (RenameResult, error) {
	if oldKey == newKey {
		return RenameResult{OldKey: oldKey, NewKey: newKey, Found: true}, nil
	}

	entries, err := LoadFile(path)
	if err != nil {
		return RenameResult{}, fmt.Errorf("rename: load %s: %w", path, err)
	}

	m := ToMap(entries)
	val, ok := m[oldKey]
	if !ok {
		return RenameResult{OldKey: oldKey, NewKey: newKey, Found: false}, nil
	}

	delete(m, oldKey)
	m[newKey] = val

	if err := WriteFile(path, m, 0600); err != nil {
		return RenameResult{}, fmt.Errorf("rename: write %s: %w", path, err)
	}

	return RenameResult{OldKey: oldKey, NewKey: newKey, Found: true}, nil
}
