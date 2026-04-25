package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForPin(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvForPin: %v", err)
	}
	return p
}

func TestPinCmd_DryRun(t *testing.T) {
	f := writeTempEnvForPin(t, "FOO=bar\nBAZ=qux\n")

	var out strings.Builder
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	err := runPin(f, nil, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// File should be unchanged
	raw, _ := os.ReadFile(f)
	if strings.Contains(string(raw), "# pinned") {
		t.Error("dry-run should not modify file")
	}
}

func TestPinCmd_PinsAllKeys(t *testing.T) {
	f := writeTempEnvForPin(t, "FOO=bar\nBAZ=qux\n")

	err := runPin(f, nil, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, _ := os.ReadFile(f)
	content := string(raw)
	if !strings.Contains(content, "# pinned") {
		t.Errorf("expected pinned markers in file, got:\n%s", content)
	}
}

func TestPinCmd_SelectedKey(t *testing.T) {
	f := writeTempEnvForPin(t, "FOO=bar\nBAZ=qux\n")

	err := runPin(f, []string{"FOO"}, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, _ := os.ReadFile(f)
	content := string(raw)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "BAZ=") && strings.Contains(line, "# pinned") {
			t.Error("BAZ should not be pinned")
		}
	}
}

func TestPinCmd_MissingFile(t *testing.T) {
	err := runPin("/nonexistent/.env", nil, "", false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
