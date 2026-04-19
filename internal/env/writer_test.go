package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFile_CreatesSortedOutput(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	data := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MIDDLE": "mid",
	}

	if err := WriteFile(p, data, 0600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(p)
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected first line ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line ZEBRA, got %s", lines[2])
	}
}

func TestWriteFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	if err := WriteFile(p, map[string]string{"K": "v"}, 0600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestFormatEnvLine_QuotesSpaces(t *testing.T) {
	line := formatEnvLine("KEY", "hello world")
	if line != `KEY="hello world"` {
		t.Errorf("unexpected line: %s", line)
	}
}

func TestFormatEnvLine_PlainValue(t *testing.T) {
	line := formatEnvLine("KEY", "simple")
	if line != "KEY=simple" {
		t.Errorf("unexpected line: %s", line)
	}
}

func TestFormatEnvLine_QuotesDollar(t *testing.T) {
	line := formatEnvLine("KEY", "$secret")
	if line != `KEY="$secret"` {
		t.Errorf("unexpected line: %s", line)
	}
}

func TestNeedsQuotes_Empty(t *testing.T) {
	if needsQuotes("") {
		t.Error("empty string should not need quotes")
	}
}
