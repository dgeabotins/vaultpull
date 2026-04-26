package env

import (
	"fmt"
	"strings"
)

// RotateOptions controls how key rotation behaves.
type RotateOptions struct {
	// Keys is the explicit list of keys to rotate; empty means all.
	Keys []string
	// Suffix is appended to the old key name to preserve the previous value.
	// Defaults to "_PREVIOUS" when empty.
	Suffix string
	// DryRun reports what would change without modifying the map.
	DryRun bool
}

// RotateResult summarises what happened during rotation.
type RotateResult struct {
	Rotated  []string // keys whose values were rotated
	Skipped  []string // keys not found in incoming
	DryRun   bool
}

// Summary returns a human-readable description of the result.
func (r RotateResult) Summary() string {
	if r.DryRun {
		return fmt.Sprintf("[dry-run] would rotate %d key(s), skip %d", len(r.Rotated), len(r.Skipped))
	}
	return fmt.Sprintf("rotated %d key(s), skipped %d", len(r.Rotated), len(r.Skipped))
}

// Rotate moves current values to <KEY><suffix> and writes incoming values into
// current. The existing map is modified in place.
func Rotate(current, incoming map[string]string, opts RotateOptions) (RotateResult, error) {
	suffix := opts.Suffix
	if suffix == "" {
		suffix = "_PREVIOUS"
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range incoming {
			keys = append(keys, k)
		}
	}

	var result RotateResult
	result.DryRun = opts.DryRun

	for _, k := range keys {
		newVal, ok := incoming[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		prevKey := strings.ToUpper(k) + suffix
		if !opts.DryRun {
			if old, exists := current[k]; exists {
				current[prevKey] = old
			}
			current[k] = newVal
		}
		result.Rotated = append(result.Rotated, k)
	}

	return result, nil
}
