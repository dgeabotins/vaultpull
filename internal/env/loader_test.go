package env

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

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	entries, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestLoadFile_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=val\n")
	entries, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Key != "KEY" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}

func TestLoadFile_StripQuotes(t *testing.T) {
	path := writeTempEnv(t, `SECRET="hello world"` + "\n")
	entries, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].Value != "hello world" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestToMap_Conversion(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	m := ToMap(entries)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}
