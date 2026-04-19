package env

import "strings"

// FilterOptions controls which keys are included.
type FilterOptions struct {
	Prefix  string
	Keys    []string
	Exclude []string
}

// Filter returns a subset of the given map based on FilterOptions.
// If Prefix is set, only keys with that prefix are included.
// If Keys is set, only those exact keys are included.
// Exclude removes keys from the result after other filters are applied.
func Filter(src map[string]string, opts FilterOptions) map[string]string {
	excludeSet := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = true
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	out := make(map[string]string)
	for k, v := range src {
		if excludeSet[k] {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if len(keySet) > 0 && !keySet[k] {
			continue
		}
		out[k] = v
	}
	return out
}
