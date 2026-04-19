package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForDedup(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDedupCmd_NoDuplicates(t *testing.T) {
	p := writeTempEnvForDedup(t, "FOO=bar\nBAZ=qux\n")
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"dedup", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no duplicates") {
		t.Errorf("expected no-duplicates message, got: %s", buf.String())
	}
}

func TestDedupCmd_WithDuplicates(t *testing.T) {
	p := writeTempEnvForDedup(t, "FOO=first\nBAR=keep\nFOO=second\n")
	buf := &bytes.Buffer{}
	_ = buf
	// build isolated command to avoid state pollution
	cmd := &cobra.Command{Use: "dedup", Args: cobra.ExactArgs(1), RunE: runDedup}
	cmd.Flags().BoolVarP(&dedupWrite, "write", "w", false, "write back")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{p})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "1 duplicate(s) removed") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDedupCmd_WriteFlag(t *testing.T) {
	p := writeTempEnvForDedup(t, "KEY=one\nKEY=two\n")
	cmd := &cobra.Command{Use: "dedup", Args: cobra.ExactArgs(1), RunE: runDedup}
	cmd.Flags().BoolVarP(&dedupWrite, "write", "w", false, "write back")
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--write", p})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(p)
	if strings.Count(string(data), "KEY=") != 1 {
		t.Errorf("expected single KEY after dedup, got:\n%s", string(data))
	}
}
