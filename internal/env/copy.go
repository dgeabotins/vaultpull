package env

import "fmt"

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	Source string
	Dest   string
	Keys   []string
}

func (r CopyResult) Summary() string {
	if len(r.Keys) == 0 {
		return fmt.Sprintf("no keys copied from %s to %s", r.Source, r.Dest)
	}
	return fmt.Sprintf("copied %d key(s) from %s to %s", len(r.Keys), r.Source, r.Dest)
}

// Copy reads selected keys from src and merges them into dst.
// If keys is empty, all keys are copied.
func Copy(src, dst string, keys []string) (CopyResult, error) {
	srcMap, err := LoadFile(src)
	if err != nil {
		return CopyResult{}, fmt.Errorf("reading source: %w", err)
	}

	selected := srcMap
	if len(keys) > 0 {
		selected = make(map[string]string, len(keys))
		for _, k := range keys {
			v, ok := srcMap[k]
			if !ok {
				return CopyResult{}, fmt.Errorf("key %q not found in %s", k, src)
			}
			selected[k] = v
		}
	}

	_, err = MergeIntoFile(dst, selected)
	if err != nil {
		return CopyResult{}, fmt.Errorf("writing dest: %w", err)
	}

	copied := make([]string, 0, len(selected))
	for k := range selected {
		copied = append(copied, k)
	}

	return CopyResult{Source: src, Dest: dst, Keys: copied}, nil
}
