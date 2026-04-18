package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForDiff(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestDiffCmd_ShowsChanges(t *testing.T) {
	old := writeTempEnvForDiff(t, "FOO=bar\nBAZ=qux\n")
	new_ := writeTempEnvForDiff(t, "FOO=changed\nNEW=value\n")

	buf := &strings.Builder{}
	cmd := &cobra.Command{}
	cmd.SetOut(buf)

	err := runDiff(cmd, old, new_)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in diff output, got: %s", out)
	}
	if !strings.Contains(out, "NEW") {
		t.Errorf("expected NEW in diff output, got: %s", out)
	}
}

func TestDiffCmd_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	old := filepath.Join(tmpDir, "old.env")
	new_ := filepath.Join(tmpDir, "new.env")

	// only create old
	os.WriteFile(old, []byte("FOO=bar\n"), 0600)

	cmd := &cobra.Command{}
	cmd.SetOut(&strings.Builder{})

	err := runDiff(cmd, old, new_)
	if err == nil {
		t.Error("expected error for missing new file")
	}
}
