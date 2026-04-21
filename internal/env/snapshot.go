package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's contents.
type Snapshot struct {
	File      string            `json:"file"`
	CapturedAt time.Time        `json:"captured_at"`
	Entries   map[string]string `json:"entries"`
}

// SnapshotResult holds the outcome of a snapshot operation.
type SnapshotResult struct {
	Snapshot *Snapshot
	Path     string
}

// TakeSnapshot reads the given env file and writes a JSON snapshot to dest.
func TakeSnapshot(src, dest string) (*SnapshotResult, error) {
	entries, err := LoadFile(src)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", src, err)
	}

	snap := &Snapshot{
		File:       src,
		CapturedAt: time.Now().UTC(),
		Entries:    entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return nil, fmt.Errorf("snapshot: write %s: %w", dest, err)
	}

	return &SnapshotResult{Snapshot: snap, Path: dest}, nil
}

// LoadSnapshot reads a previously saved snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: parse %s: %w", path, err)
	}

	return &snap, nil
}
