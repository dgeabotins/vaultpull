package audit

// NoopLogger satisfies the same interface but discards all entries.
// Useful when audit logging is not configured.
type NoopLogger struct{}

// Record discards the entry and returns nil.
func (n *NoopLogger) Record(_ Entry) error {
	return nil
}
