package env

import (
	"fmt"
	"strings"
)

// DuplicateEntry represents a key found more than once in an env file.
type DuplicateEntry struct {
	Key   string
	Lines []int
}

// DedupResult holds the outcome of a deduplication pass.
type DedupResult struct {
	Duplicates []DuplicateEntry
	Removed    int
}

// Summary returns a human-readable summary of the dedup result.
func (r DedupResult) Summary() string {
	if len(r.Duplicates) == 0 {
		return "no duplicates found"
	}
	var sb strings.Builder
	for _, d := range r.Duplicates {
		fmt.Fprintf(&sb, "key %q duplicated on lines %v\n", d.Key, d.Lines)
	}
	fmt.Fprintf(&sb, "%d duplicate(s) removed", r.Removed)
	return strings.TrimSpace(sb.String())
}

// Dedup reads an env map produced by LoadFile and returns a deduplicated map
// along with a DedupResult describing what was removed. The first occurrence
// of each key is kept.
func Dedup(lines []string) ([]string, DedupResult) {
	seen := map[string]int{}
	dupeMap := map[string]*DuplicateEntry{}
	var out []string
	removed := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			out = append(out, line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		lineNum := i + 1
		if firstLine, exists := seen[key]; exists {
			if dupeMap[key] == nil {
				dupeMap[key] = &DuplicateEntry{Key: key, Lines: []int{firstLine}}
			}
			dupeMap[key].Lines = append(dupeMap[key].Lines, lineNum)
			removed++
			continue
		}
		seen[key] = lineNum
		out = append(out, line)
	}

	var dupes []DuplicateEntry
	for _, d := range dupeMap {
		dupes = append(dupes, *d)
	}
	return out, DedupResult{Duplicates: dupes, Removed: removed}
}
