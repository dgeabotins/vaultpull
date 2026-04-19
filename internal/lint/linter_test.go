package lint

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
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

func TestCheck_NoIssues(t *testing.T) {
	path := writeTempEnv(t, "APP_KEY=value\nDB_HOST=localhost\n")
	res, err := Check(path)
	if err != nil {
		t.Fatal(err)
	}
	if res.HasIssues() {
		t.Errorf("expected no issues, got: %s", res.Summary())
	}
}

func TestCheck_InvalidKey(t *testing.T) {
	path := writeTempEnv(t, "lower_key=value\n")
	res, err := Check(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasIssues() {
		t.Error("expected issue for lowercase key")
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_KEY=\n")
	res, err := Check(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasIssues() {
		t.Error("expected issue for empty value")
	}
}

func TestCheck_MissingSeparator(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	res, err := Check(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasIssues() {
		t.Error("expected issue for missing separator")
	}
}

func TestCheck_FileNotFound(t *testing.T) {
	_, err := Check("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSummary_NoIssues(t *testing.T) {
	r := &Result{}
	if r.Summary() != "no issues found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
