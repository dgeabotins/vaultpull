package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Profile   string    `json:"profile"`
	Path      string    `json:"path"`
	Keys      []string  `json:"keys"`
	OutputFile string   `json:"output_file"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// New creates a new Logger that appends to the given file path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Record writes an audit entry to the log file.
func (l *Logger) Record(e Entry) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}
