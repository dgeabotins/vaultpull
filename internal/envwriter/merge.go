package envwriter

import (
	"bufio"
	"os"
	"strings"
)

// Merge reads existing key=value pairs from path and overlays newSecrets on
// top, preserving any keys not present in newSecrets. Returns the merged map.
func Merge(path string, newSecrets map[string]string) (map[string]string, error) {
	existing, err := readEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	for k, v := range newSecrets {
		existing[k] = v
	}
	return existing, nil
}

// readEnvFile parses a .env file into a map, ignoring comments and blank lines.
func readEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return map[string]string{}, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = strings.Trim(parts[1], `"`)
	}
	return result, scanner.Err()
}
