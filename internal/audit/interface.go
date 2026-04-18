package audit

// Recorder is the interface implemented by Logger and NoopLogger.
type Recorder interface {
	Record(Entry) error
}

// NewRecorder returns a real Logger when path is non-empty,
// otherwise returns a NoopLogger.
func NewRecorder(path string) Recorder {
	if path == "" {
		return &NoopLogger{}
	}
	return New(path)
}
