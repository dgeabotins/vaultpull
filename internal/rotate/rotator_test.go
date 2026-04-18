package rotate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackup_CreatesBackupFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, ".env")
	if err := os.WriteFile(src, []byte("KEY=val\n"), 0600); err != nil {
		t.Fatal(err)
	}

	r := New(filepath.Join(tmp, "backups"), 5)
	dest, err := r.Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest == "" {
		t.Fatal("expected non-empty backup path")
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("backup file not readable: %v", err)
	}
	if string(data) != "KEY=val\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestBackup_NonExistentSrc(t *testing.T) {
	tmp := t.TempDir()
	r := New(filepath.Join(tmp, "backups"), 5)
	dest, err := r.Backup(filepath.Join(tmp, "missing.env"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if dest != "" {
		t.Errorf("expected empty dest for missing src, got %q", dest)
	}
}

func TestBackup_PrunesOldBackups(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, ".env")
	backupDir := filepath.Join(tmp, "backups")
	if err := os.WriteFile(src, []byte("A=1"), 0600); err != nil {
		t.Fatal(err)
	}

	r := New(backupDir, 3)
	for i := 0; i < 5; i++ {
		if _, err := r.Backup(src); err != nil {
			t.Fatalf("backup %d failed: %v", i, err)
		}
	}

	matches, err := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) > 3 {
		t.Errorf("expected at most 3 backups, got %d", len(matches))
	}
}

func TestNew_DefaultMaxBackups(t *testing.T) {
	r := New("/tmp", 0)
	if r.maxBackups != 5 {
		t.Errorf("expected default maxBackups=5, got %d", r.maxBackups)
	}
}
