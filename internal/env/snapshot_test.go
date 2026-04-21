package env

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForSnapshot(t *testing.T, content string) string {
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

func TestTakeSnapshot_CreatesFile(t *testing.T) {
	src := writeTempEnvForSnapshot(t, "FOO=bar\nBAZ=qux\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	res, err := TakeSnapshot(src, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Path != dest {
		t.Errorf("expected path %s, got %s", dest, res.Path)
	}
	if res.Snapshot.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", res.Snapshot.Entries["FOO"])
	}
	if res.Snapshot.Entries["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", res.Snapshot.Entries["BAZ"])
	}
}

func TestTakeSnapshot_FilePermissions(t *testing.T) {
	src := writeTempEnvForSnapshot(t, "KEY=value\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	if _, err := TakeSnapshot(src, dest); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestLoadSnapshot_RoundTrip(t *testing.T) {
	src := writeTempEnvForSnapshot(t, "A=1\nB=2\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	if _, err := TakeSnapshot(src, dest); err != nil {
		t.Fatal(err)
	}

	snap, err := LoadSnapshot(dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Entries["A"] != "1" {
		t.Errorf("expected A=1, got %s", snap.Entries["A"])
	}
	if snap.File != src {
		t.Errorf("expected file %s, got %s", src, snap.File)
	}
}

func TestTakeSnapshot_MissingSource(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "snap.json")
	_, err := TakeSnapshot("/nonexistent/.env", dest)
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("not json")
	f.Close()

	_, err = LoadSnapshot(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestTakeSnapshot_ContainsCapturedAt(t *testing.T) {
	src := writeTempEnvForSnapshot(t, "X=y\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	if _, err := TakeSnapshot(src, dest); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dest)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["captured_at"]; !ok {
		t.Error("expected captured_at field in snapshot JSON")
	}
}
