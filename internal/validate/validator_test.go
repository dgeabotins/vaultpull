package validate

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

func TestCheckEnvFile_AllPresent(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	results, err := CheckEnvFile(path, []string{"FOO", "BAZ"})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Missing || r.Empty {
			t.Errorf("expected OK for %s, got %s", r.Key, r)
		}
	}
}

func TestCheckEnvFile_MissingKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	results, err := CheckEnvFile(path, []string{"FOO", "MISSING"})
	if err != nil {
		t.Fatal(err)
	}
	if !results[1].Missing {
		t.Errorf("expected MISSING for MISSING key")
	}
}

func TestCheckEnvFile_EmptyValue(t *testing.T) {
	path := writeTempEnv(t, "FOO=\n")
	results, err := CheckEnvFile(path, []string{"FOO"})
	if err != nil {
		t.Fatal(err)
	}
	if !results[0].Empty {
		t.Errorf("expected EMPTY for FOO")
	}
}

func TestCheckEnvFile_FileNotFound(t *testing.T) {
	_, err := CheckEnvFile("/nonexistent/.env", []string{"X"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestHasFailures_True(t *testing.T) {
	results := []Result{{Key: "A", Missing: true}}
	if !HasFailures(results) {
		t.Error("expected HasFailures to return true")
	}
}

func TestHasFailures_False(t *testing.T) {
	results := []Result{{Key: "A"}, {Key: "B"}}
	if HasFailures(results) {
		t.Error("expected HasFailures to return false")
	}
}

func TestResult_String(t *testing.T) {
	if got := (Result{Key: "X", Missing: true}).String(); got != "MISSING  X" {
		t.Errorf("unexpected: %s", got)
	}
	if got := (Result{Key: "Y", Empty: true}).String(); got != "EMPTY    Y" {
		t.Errorf("unexpected: %s", got)
	}
	if got := (Result{Key: "Z"}).String(); got != "OK       Z" {
		t.Errorf("unexpected: %s", got)
	}
}
