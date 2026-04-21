package env

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging env maps.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing keys and only adds new ones.
	StrategyKeepExisting
	// StrategyError returns an error on any key conflict.
	StrategyError
)

// MergeOptions configures the behaviour of MergeMap.
type MergeOptions struct {
	Strategy MergeStrategy
	// Prefix, if set, only merges keys that start with this prefix.
	Prefix string
}

// MergeResult holds statistics from a MergeMap call.
type MergeResult struct {
	Added     int
	Updated   int
	Skipped   int
	Conflicts []string
}

// MergeMap merges src into dst according to opts.
// dst is modified in place. Returns a MergeResult and any error.
func MergeMap(dst, src map[string]string, opts MergeOptions) (MergeResult, error) {
	var result MergeResult

	for k, v := range src {
		if opts.Prefix != "" {
			if len(k) < len(opts.Prefix) || k[:len(opts.Prefix)] != opts.Prefix {
				continue
			}
		}

		existing, exists := dst[k]
		if !exists {
			dst[k] = v
			result.Added++
			continue
		}

		if existing == v {
			result.Skipped++
			continue
		}

		switch opts.Strategy {
		case StrategyOverwrite:
			dst[k] = v
			result.Updated++
		case StrategyKeepExisting:
			result.Skipped++
		case StrategyError:
			result.Conflicts = append(result.Conflicts, k)
		}
	}

	if len(result.Conflicts) > 0 {
		return result, fmt.Errorf("merge conflict on keys: %v", result.Conflicts)
	}
	return result, nil
}
