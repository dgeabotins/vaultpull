package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// WriteFile writes a map of key-value pairs to a file in .env format.
// Keys are sorted alphabetically for deterministic output.
func WriteFile(path string, data map[string]string, perm os.FileMode) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := data[k]
		sb.WriteString(formatEnvLine(k, v))
		sb.WriteByte('\n')
	}

	return os.WriteFile(path, []byte(sb.String()), perm)
}

// formatEnvLine returns a single KEY=VALUE line, quoting the value if needed.
func formatEnvLine(key, value string) string {
	if needsQuotes(value) {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return fmt.Sprintf(`%s="%s"`, key, escaped)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// needsQuotes reports whether a value should be wrapped in double quotes.
func needsQuotes(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch == ' ' || ch == '\t' || ch == '#' || ch == '$' || ch == '\\' {
			return true
		}
	}
	return false
}
