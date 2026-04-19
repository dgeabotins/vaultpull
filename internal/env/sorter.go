package env

import (
	"sort"
	"strings"
)

// SortOptions controls how env keys are sorted.
type SortOptions struct {
	CaseInsensitive bool
	Reverse         bool
	PrefixFirst     []string // prefixes to sort before others
}

// SortKeys returns a sorted copy of the provided keys slice.
func SortKeys(keys []string, opts SortOptions) []string {
	out := make([]string, len(keys))
	copy(out, keys)

	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]

		// Priority prefix ordering
		for _, p := range opts.PrefixFirst {
			aHas := strings.HasPrefix(a, p)
			bHas := strings.HasPrefix(b, p)
			if aHas && !bHas {
				return true
			}
			if !aHas && bHas {
				return false
			}
		}

		if opts.CaseInsensitive {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if opts.Reverse {
			return a > b
		}
		return a < b
	})

	return out
}

// SortMap returns keys of the map in sorted order using SortKeys.
func SortMap(m map[string]string, opts SortOptions) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return SortKeys(keys, opts)
}
