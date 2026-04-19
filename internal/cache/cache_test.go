package cache

import (
	"os"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "vaultpull-cache-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestSet_And_Get(t *testing.T) {
	c := New(tempDir(t), 5*time.Minute)
	secrets := map[string]string{"KEY": "value", "TOKEN": "abc123"}

	if err := c.Set("dev", secrets); err != nil {
		t.Fatalf("Set: %v", err)
	}

	e, ok := c.Get("dev")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if e.Secrets["KEY"] != "value" {
		t.Errorf("got %q, want %q", e.Secrets["KEY"], "value")
	}
}

func TestGet_Miss_NoFile(t *testing.T) {
	c := New(tempDir(t), 5*time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestGet_Miss_Expired(t *testing.T) {
	c := New(tempDir(t), 1*time.Millisecond)
	if err := c.Set("prod", map[string]string{"A": "1"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("prod")
	if ok {
		t.Fatal("expected expired cache miss")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := New(tempDir(t), 5*time.Minute)
	_ = c.Set("staging", map[string]string{"X": "y"})

	if err := c.Invalidate("staging"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	_, ok := c.Get("staging")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestInvalidate_NonExistent(t *testing.T) {
	c := New(tempDir(t), 5*time.Minute)
	if err := c.Invalidate("ghost"); err != nil {
		t.Errorf("expected no error for missing entry, got %v", err)
	}
}
