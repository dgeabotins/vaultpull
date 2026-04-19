package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForCopyCmd(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0600)
	return f.Name()
}

func TestCopyCmd_AllKeys(t *testing.T) {
	src := writeTempEnvForCopyCmd(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "copied") {
		t.Errorf("expected 'copied' in output, got: %s", buf.String())
	}
}

func TestCopyCmd_SelectedKeys(t *testing.T) {
	src := writeTempEnvForCopyCmd(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"copy", "--keys", "FOO", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "1 key") {
		t.Errorf("expected '1 key' in output, got: %s", buf.String())
	}
}

func TestCopyCmd_MissingSource(t *testing.T) {
	c := &cobra.Command{Use: "copy", RunE: runCopy, Args: cobra.ExactArgs(2)}
	c.Flags().StringVarP(&copyKeys, "keys", "k", "", "")
	c.SetArgs([]string{"/no/such.env", "/tmp/dst.env"})
	if err := c.Execute(); err == nil {
		t.Error("expected error for missing source")
	}
}
