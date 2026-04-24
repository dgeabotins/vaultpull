package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForCast(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestCastCmd_DryRun_ShowsChanges(t *testing.T) {
	p := writeTempEnvForCast(t, "FLAG=yes\nPORT=8080\n")

	buf := new(strings.Builder)
	rootCmd.SetOut(buf)

	// re-register output to stdout for test capture
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := executeTestCmd(rootCmd, "cast",
		"--file", p,
		"--hint", "FLAG=bool",
		"--hint", "PORT=int",
		"--strict-bool",
		"--dry-run",
	)

	w.Close()
	out := make([]byte, 4096)
	n, _ := r.Read(out)
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := string(out[:n])
	if !strings.Contains(output, "changed") && !strings.Contains(output, "no values") {
		t.Errorf("expected change summary in output, got: %s", output)
	}
}

func TestCastCmd_MissingFile(t *testing.T) {
	err := executeTestCmd(rootCmd, "cast", "--file", "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCastCmd_InvalidHint(t *testing.T) {
	p := writeTempEnvForCast(t, "X=1\n")
	err := executeTestCmd(rootCmd, "cast", "--file", p, "--hint", "BADFORMAT")
	if err == nil {
		t.Error("expected error for malformed hint")
	}
}

// executeTestCmd is a helper to run a cobra command by name.
func executeTestCmd(root *cobra.Command, args ...string) error {
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return err
}
