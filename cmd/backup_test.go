package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBackupCmd_Success(t *testing.T) {
	tmp := t.TempDir()
	envFile := filepath.Join(tmp, ".env")
	if err := os.WriteFile(envFile, []byte("SECRET=abc\n"), 0600); err != nil {
		t.Fatal(err)
	}

	bDir := filepath.Join(tmp, "backups")

	cmd := &cobra.Command{}
	_ = cmd

	// Exercise via rootCmd
	rootCmd.SetArgs([]string{
		"backup", envFile,
		"--backup-dir", bDir,
		"--max-backups", "3",
	})

	buf := new(strings.Builder)
	rootCmd.SetOut(buf)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "backed up to") {
		t.Errorf("expected success message, got: %q", buf.String())
	}

	matches, _ := filepath.Glob(filepath.Join(bDir, ".env.*.bak"))
	if len(matches) != 1 {
		t.Errorf("expected 1 backup file, got %d", len(matches))
	}
}

func TestBackupCmd_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	bDir := filepath.Join(tmp, "backups")

	rootCmd.SetArgs([]string{
		"backup", filepath.Join(tmp, "nonexistent.env"),
		"--backup-dir", bDir,
	})

	buf := new(strings.Builder)
	rootCmd.SetOut(buf)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "nothing to back up") {
		t.Errorf("expected 'nothing to back up', got: %q", buf.String())
	}
}
