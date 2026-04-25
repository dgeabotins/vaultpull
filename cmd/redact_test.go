package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForRedact(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRedactCmd_DryRun(t *testing.T) {
	file := writeTempEnvForRedact(t, "DB_PASSWORD=secret\nAPP_NAME=myapp\n")

	out := captureRedactOutput(t, []string{
		"redact", file,
		"--key", "DB_PASSWORD",
		"--dry-run",
	})

	if !strings.Contains(out, "DB_PASSWORD=***") {
		t.Errorf("expected redacted output, got: %s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME unchanged, got: %s", out)
	}
}

func TestRedactCmd_WritesFile(t *testing.T) {
	file := writeTempEnvForRedact(t, "TOKEN=abc123\nHOST=localhost\n")

	captureRedactOutput(t, []string{
		"redact", file,
		"--key", "TOKEN",
	})

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(data), "TOKEN=***") {
		t.Errorf("expected TOKEN redacted in file, got: %s", string(data))
	}
}

func TestRedactCmd_MissingFile(t *testing.T) {
	err := runRedact("/nonexistent/.env", nil, nil, "***", "", true)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func captureRedactOutput(t *testing.T, args []string) string {
	t.Helper()
	old := rootCmd
	rootCmd = &cobra.Command{Use: "vaultpull"}
	defer func() { rootCmd = old }()

	init_redact_cmd(rootCmd)

	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()
	return buf.String()
}

func init_redact_cmd(root *cobra.Command) {
	var keys []string
	var patterns []string
	var placeholder string
	var dryRun bool
	var output string

	cmd := &cobra.Command{
		Use:  "redact <file>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedact(args[0], keys, patterns, placeholder, output, dryRun)
		},
	}
	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "")
	cmd.Flags().StringSliceVarP(&patterns, "pattern", "p", nil, "")
	cmd.Flags().StringVar(&placeholder, "placeholder", "***", "")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "")
	cmd.Flags().StringVarP(&output, "output", "o", "", "")
	root.AddCommand(cmd)
}
