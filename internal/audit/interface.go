package audit

// Recorder is the interface implemented by Logger and NoopLogger.
type Recorder interface {
	Record(Entry) error
}

// NewRecorder returns a real Logger when path is non-empty,
// otherwise returns a NoopLogger. The caller is responsible for
// closing the underlying logger when it is no longer needed.
func NewRecorder(path string) Recorder {
	if path == "" {
		return &NoopLogger{}
	}
	return New(path)
}

// IsNoop reports whether r is a no-op recorder that discards all entries.
func IsNoop(r Recorder) bool {
	_, ok := r.(*NoopLogger)
	return ok
}
