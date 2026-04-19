package env

import (
	"testing"
)

func TestSortKeys_Default(t *testing.T) {
	keys := []string{"ZEBRA", "ALPHA", "MANGO"}
	got := SortKeys(keys, SortOptions{})
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range got {
		if k != want[i] {
			t.Errorf("pos %d: got %q want %q", i, k, want[i])
		}
	}
}

func TestSortKeys_Reverse(t *testing.T) {
	keys := []string{"ALPHA", "MANGO", "ZEBRA"}
	got := SortKeys(keys, SortOptions{Reverse: true})
	want := []string{"ZEBRA", "MANGO", "ALPHA"}
	for i, k := range got {
		if k != want[i] {
			t.Errorf("pos %d: got %q want %q", i, k, want[i])
		}
	}
}

func TestSortKeys_CaseInsensitive(t *testing.T) {
	keys := []string{"zebra", "ALPHA", "mango"}
	got := SortKeys(keys, SortOptions{CaseInsensitive: true})
	want := []string{"ALPHA", "mango", "zebra"}
	for i, k := range got {
		if k != want[i] {
			t.Errorf("pos %d: got %q want %q", i, k, want[i])
		}
	}
}

func TestSortKeys_PrefixFirst(t *testing.T) {
	keys := []string{"DB_HOST", "APP_NAME", "DB_PORT", "LOG_LEVEL"}
	got := SortKeys(keys, SortOptions{PrefixFirst: []string{"APP_"}})
	if got[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %q", got[0])
	}
}

func TestSortMap_ReturnsAllKeys(t *testing.T) {
	m := map[string]string{"B": "2", "A": "1", "C": "3"}
	got := SortMap(m, SortOptions{})
	if len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
	if got[0] != "A" || got[1] != "B" || got[2] != "C" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestSortKeys_OriginalUnmodified(t *testing.T) {
	orig := []string{"C", "A", "B"}
	copy := []string{"C", "A", "B"}
	SortKeys(orig, SortOptions{})
	for i, v := range orig {
		if v != copy[i] {
			t.Errorf("original slice was modified at index %d", i)
		}
	}
}
