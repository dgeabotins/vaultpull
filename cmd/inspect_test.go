package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempEnvForInspect(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestInspectCmd_ShowsEntries(t *testing.T) {
	path := writeTempEnvForInspect(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"inspect", "--file", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected localhost in output, got: %s", out)
	}
}

func TestInspectCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"inspect", "--file", "/no/such/.env"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInspectEnvFile_ReturnsEntries(t *testing.T) {
	path := writeTempEnvForInspect(t, "TOKEN=abc123\n")
	entries, err := InspectEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Key != "TOKEN" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}
