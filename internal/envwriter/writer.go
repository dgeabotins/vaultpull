package envwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	OutputPath string
}

// New creates a new Writer targeting the given output path.
func New(outputPath string) *Writer {
	return &Writer{OutputPath: outputPath}
}

// Write serializes the provided secrets map into .env format and writes
// it to the configured output path, creating parent directories as needed.
func (w *Writer) Write(secrets map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(w.OutputPath), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	var sb strings.Builder
	for k, v := range secrets {
		sb.WriteString(formatLine(k, v))
	}

	if err := os.WriteFile(w.OutputPath, []byte(sb.String()), 0o600); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}
	return nil
}

// formatLine returns a single KEY=VALUE\n line, quoting the value when needed.
func formatLine(key, value string) string {
	if strings.ContainsAny(value, " \t\n#") {
		value = `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
	}
	return fmt.Sprintf("%s=%s\n", key, value)
}
