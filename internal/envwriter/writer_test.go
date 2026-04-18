package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "subdir", ".env")

	w := New(output)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	content := string(data)
	for k, v := range secrets {
		expected := k + "=" + v
		if !strings.Contains(content, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, content)
		}
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	w := New(output)
	if err := w.Write(map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(output)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}

func TestFormatLine_QuotesSpecialValues(t *testing.T) {
	line := formatLine("MSG", "hello world")
	if !strings.Contains(line, `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", line)
	}
}

func TestFormatLine_PlainValue(t *testing.T) {
	line := formatLine("TOKEN", "abc123")
	expected := "TOKEN=abc123\n"
	if line != expected {
		t.Errorf("expected %q, got %q", expected, line)
	}
}
