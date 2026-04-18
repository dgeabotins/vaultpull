package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForValidate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnvForValidate: %v", err)
	}
	return p
}

func TestValidateCmd_AllPresent(t *testing.T) {
	envFile := writeTempEnvForValidate(t, "FOO=bar\nBAZ=qux\n")

	buf := new(strings.Builder)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"validate", envFile, "--keys", "FOO,BAZ"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCmd_MissingKey(t *testing.T) {
	envFile := writeTempEnvForValidate(t, "FOO=bar\n")

	cmd := &cobra.Command{}
	_ = cmd

	err := runValidate(envFile, []string{"FOO", "MISSING"}, false)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateCmd_WarnOnly(t *testing.T) {
	envFile := writeTempEnvForValidate(t, "FOO=bar\n")

	err := runValidate(envFile, []string{"FOO", "MISSING"}, true)
	if err != nil {
		t.Fatalf("expected no error in warn-only mode, got: %v", err)
	}
}

func TestValidateCmd_FileNotFound(t *testing.T) {
	err := runValidate("/nonexistent/.env", []string{"FOO"}, false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
