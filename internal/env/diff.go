package env

import "fmt"

// DiffResult holds the comparison between two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
	Same    map[string]string
}

// Diff compares two env maps and returns a DiffResult.
func Diff(old, next map[string]string) DiffResult {
	r := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
		Same:    make(map[string]string),
	}
	for k, v := range next {
		ov, ok := old[k]
		if !ok {
			r.Added[k] = v
		} else if ov != v {
			r.Changed[k] = [2]string{ov, v}
		} else {
			r.Same[k] = v
		}
	}
	for k, v := range old {
		if _, ok := next[k]; !ok {
			r.Removed[k] = v
		}
	}
	return r
}

// HasChanges returns true if there are any added, removed, or changed keys.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary string.
func (d DiffResult) Summary() string {
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed, %d unchanged",
		len(d.Added), len(d.Removed), len(d.Changed), len(d.Same))
}
