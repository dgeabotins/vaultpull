package diff

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the diff outcome for a single key.
type Result struct {
	Key    string
	OldVal string
	NewVal string
	Status string // added | removed | changed | unchanged
}

// Compare reads two env files and returns per-key diff results.
func Compare(oldPath, newPath string) ([]Result, error) {
	oldMap, err := parseEnvFile(oldPath)
	if err != nil {
		return nil, fmt.Errorf("reading old file: %w", err)
	}
	newMap, err := parseEnvFile(newPath)
	if err != nil {
		return nil, fmt.Errorf("reading new file: %w", err)
	}

	seen := map[string]bool{}
	var results []Result

	for k, ov := range oldMap {
		seen[k] = true
		if nv, ok := newMap[k]; !ok {
			results = append(results, Result{Key: k, OldVal: ov, Status: "removed"})
		} else if ov != nv {
			results = append(results, Result{Key: k, OldVal: ov, NewVal: nv, Status: "changed"})
		} else {
			results = append(results, Result{Key: k, OldVal: ov, NewVal: nv, Status: "unchanged"})
		}
	}

	for k, nv := range newMap {
		if !seen[k] {
			results = append(results, Result{Key: k, NewVal: nv, Status: "added"})
		}
	}

	return results, nil
}

// Summary returns a human-readable summary line.
func Summary(results []Result) string {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Status]++
	}
	return fmt.Sprintf("+%d added  ~%d changed  -%d removed  =%d unchanged",
		counts["added"], counts["changed"], counts["removed"], counts["unchanged"])
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m, scanner.Err()
}
