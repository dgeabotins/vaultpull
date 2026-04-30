package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForLint(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func executeLintCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"lint"}, args...))
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func TestLintCmd_NoIssues(t *testing.T) {
	p := writeTempEnvForLint(t, "DATABASE_URL=postgres://localhost/db\nAPI_KEY=abc123\n")
	out, err := executeLintCmd(t, []string{"--file", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no lint issues") {
		t.Errorf("expected clean output, got: %s", out)
	}
}

func TestLintCmd_InvalidKey(t *testing.T) {
	p := writeTempEnvForLint(t, "bad-key=value\n")
	// Use a sub-command instance to avoid os.Exit in tests
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "lint", RunE: runLint}
	cmd.Flags().StringVarP(&lintFile, "file", "f", ".env", "")
	cmd.Flags().BoolVar(&lintWarnOnly, "warn-only", false, "")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_ = cmd.ParseFlags([]string{"--file", p, "--warn-only"})
	_ = cmd.RunE(cmd, nil)
	if !strings.Contains(buf.String(), "bad-key") {
		t.Errorf("expected key in output, got: %s", buf.String())
	}
}

func TestLintCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{Use: "lint", RunE: runLint}
	cmd.Flags().StringVarP(&lintFile, "file", "f", ".env", "")
	cmd.Flags().BoolVar(&lintWarnOnly, "warn-only", false, "")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	_ = cmd.ParseFlags([]string{"--file", "/nonexistent/.env"})
	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
