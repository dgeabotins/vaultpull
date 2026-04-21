package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForPromote(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestPromoteCmd_AllKeys(t *testing.T) {
	src := writeTempEnvForPromote(t, "A=1\nB=2\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	rootCmd.SetArgs([]string{"promote", "--src", src, "--dst", dst})
	out := captureOutput(t, rootCmd)

	if !strings.Contains(out, "promoted: 2") {
		t.Errorf("expected promoted count in output, got: %s", out)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("dst file not created: %v", err)
	}
	if !strings.Contains(string(data), "A=") {
		t.Error("expected A in dst file")
	}
}

func TestPromoteCmd_DryRun(t *testing.T) {
	src := writeTempEnvForPromote(t, "X=hello\n")
	dst := filepath.Join(t.TempDir(), "dst.env")

	rootCmd.SetArgs([]string{"promote", "--src", src, "--dst", dst, "--dry-run"})
	out := captureOutput(t, rootCmd)

	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run notice, got: %s", out)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Error("dst file should not exist after dry-run")
	}
}

func TestPromoteCmd_MissingSource(t *testing.T) {
	rootCmd.SetArgs([]string{"promote", "--src", "/no/such/file.env", "--dst", "/tmp/out.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

// captureOutput executes cmd and returns combined stdout as string.
func captureOutput(t *testing.T, cmd *cobra.Command) string {
	t.Helper()
	var sb strings.Builder
	cmd.SetOut(&sb)
	cmd.SetErr(&sb)
	_ = cmd.Execute()
	return sb.String()
}
