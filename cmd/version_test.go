package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/version"
)

func TestVersionCmd_Output(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, version.Version) {
		t.Errorf("version command output %q does not contain version %q", out, version.Version)
	}
}
