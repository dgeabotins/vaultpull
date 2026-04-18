package audit

import (
	"path/filepath"
	"testing"
)

func TestNewRecorder_ReturnsNoop(t *testing.T) {
	r := NewRecorder("")
	if _, ok := r.(*NoopLogger); !ok {
		t.Errorf("expected *NoopLogger, got %T", r)
	}
	// Should not error
	if err := r.Record(Entry{Status: "ok"}); err != nil {
		t.Errorf("noop record: %v", err)
	}
}

func TestNewRecorder_ReturnsLogger(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	r := NewRecorder(path)
	if _, ok := r.(*Logger); !ok {
		t.Errorf("expected *Logger, got %T", r)
	}
}
