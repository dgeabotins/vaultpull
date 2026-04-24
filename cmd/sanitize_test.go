package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForSanitize(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnvForSanitize: %v", err)
	}
	return p
}

func TestSanitizeCmd_NoChanges(t *testing.T) {
	p := writeTempEnvForSanitize(t, "KEY=clean\nFOO=bar\n")

	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"sanitize", "--trim", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes needed") {
		t.Errorf("expected 'no changes needed', got: %q", buf.String())
	}
}

func TestSanitizeCmd_ShowsChanges(t *testing.T) {
	p := writeTempEnvForSanitize(t, "KEY=  hello  \nFOO=bar\n")

	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"sanitize", "--trim", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY") {
		t.Errorf("expected KEY in output, got: %q", buf.String())
	}
}

func TestSanitizeCmd_DryRun(t *testing.T) {
	p := writeTempEnvForSanitize(t, "KEY=  spaced  \n")
	original, _ := os.ReadFile(p)

	buf := &strings.Builder{}
	// reset flags between tests
	var cmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "sanitize <file>" {
			cmd = c
			break
		}
	}
	if cmd != nil {
		cmd.ResetFlags()
	}

	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"sanitize", "--trim", "--dry-run", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after, _ := os.ReadFile(p)
	if string(after) != string(original) {
		t.Errorf("dry-run should not modify file")
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("expected dry-run notice in output")
	}
}

func TestSanitizeCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"sanitize", "/nonexistent/.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
