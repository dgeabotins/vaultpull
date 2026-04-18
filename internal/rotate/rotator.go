package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator handles backup and rotation of existing .env files.
type Rotator struct {
	backupDir string
	maxBackups int
}

// New creates a new Rotator. backupDir is where backups are stored.
func New(backupDir string, maxBackups int) *Rotator {
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Rotator{backupDir: backupDir, maxBackups: maxBackups}
}

// Backup copies src to the backup directory with a timestamp suffix.
// Returns the backup path or an error.
func (r *Rotator) Backup(src string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // nothing to back up
		}
		return "", fmt.Errorf("read source: %w", err)
	}

	if err := os.MkdirAll(r.backupDir, 0700); err != nil {
		return "", fmt.Errorf("create backup dir: %w", err)
	}

	base := filepath.Base(src)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dest := filepath.Join(r.backupDir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}

	if err := r.pruneOld(base); err != nil {
		return dest, fmt.Errorf("prune old backups: %w", err)
	}

	return dest, nil
}

// pruneOld removes oldest backups for the given base filename if count exceeds maxBackups.
func (r *Rotator) pruneOld(base string) error {
	pattern := filepath.Join(r.backupDir, fmt.Sprintf("%s.*.bak", base))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for len(matches) > r.maxBackups {
		oldest := matches[0]
		if err := os.Remove(oldest); err != nil {
			return err
		}
		matches = matches[1:]
	}
	return nil
}
