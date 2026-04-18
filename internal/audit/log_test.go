package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRecord_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	l := New(logPath)
	e := Entry{
		Timestamp:  time.Now().UTC(),
		Profile:    "staging",
		Path:       "secret/app",
		Keys:       []string{"DB_URL", "API_KEY"},
		OutputFile: ".env",
		Status:     "success",
	}

	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("could not open log: %v", err)
	}
	defer f.Close()

	var got Entry
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected one log line")
	}
	if err := json.Unmarshal(scanner.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Profile != "staging" {
		t.Errorf("profile: got %q, want %q", got.Profile, "staging")
	}
	if got.Status != "success" {
		t.Errorf("status: got %q, want %q", got.Status, "success")
	}
	if len(got.Keys) != 2 {
		t.Errorf("keys length: got %d, want 2", len(got.Keys))
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	dir := t.TempDir()
	l := New(filepath.Join(dir, "audit.log"))

	e := Entry{Profile: "prod", Status: "success"}
	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	l := New(filepath.Join(dir, "audit.log"))

	for i := 0; i < 3; i++ {
		if err := l.Record(Entry{Status: "success"}); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}

	data, _ := os.ReadFile(filepath.Join(dir, "audit.log"))
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	if lines != 3 {
		t.Errorf("expected 3 lines, got %d", lines)
	}
}
