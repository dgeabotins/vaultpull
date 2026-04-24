package env

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the string used to join key segments (default: "_").
	Separator string
	// Uppercase converts all keys to uppercase.
	Uppercase bool
	// Prefix is prepended to every resulting key.
	Prefix string
}

// FlattenResult holds the output of a Flatten operation.
type FlattenResult struct {
	Flattened map[string]string
	Original  map[string]string
	Renamed   int
}

// Summary returns a human-readable description of the flatten result.
func (r FlattenResult) Summary() string {
	if r.Renamed == 0 {
		return fmt.Sprintf("%d keys processed, no renames needed", len(r.Flattened))
	}
	return fmt.Sprintf("%d keys processed, %d renamed", len(r.Flattened), r.Renamed)
}

// Flatten normalises env map keys by replacing any dots or dashes in key
// names with the configured separator, optionally uppercasing and prefixing.
// It is non-destructive: the original map is left unchanged.
func Flatten(src map[string]string, opts FlattenOptions) FlattenResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	out := make(map[string]string, len(src))
	renamed := 0

	// Process keys in deterministic order so tests are stable.
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		newKey := k
		newKey = strings.ReplaceAll(newKey, ".", opts.Separator)
		newKey = strings.ReplaceAll(newKey, "-", opts.Separator)
		if opts.Uppercase {
			newKey = strings.ToUpper(newKey)
		}
		if opts.Prefix != "" {
			newKey = opts.Prefix + newKey
		}
		if newKey != k {
			renamed++
		}
		out[newKey] = src[k]
	}

	return FlattenResult{
		Flattened: out,
		Original:  src,
		Renamed:   renamed,
	}
}
