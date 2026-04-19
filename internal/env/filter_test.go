package env

import (
	"testing"
)

func TestFilter_ByPrefix(t *testing.T) {
	src := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	got := Filter(src, FilterOptions{Prefix: "APP_"})
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if got["APP_HOST"] != "localhost" || got["APP_PORT"] != "8080" {
		t.Error("unexpected values")
	}
}

func TestFilter_ByKeys(t *testing.T) {
	src := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	got := Filter(src, FilterOptions{Keys: []string{"A", "C"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["B"]; ok {
		t.Error("B should be excluded")
	}
}

func TestFilter_Exclude(t *testing.T) {
	src := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	got := Filter(src, FilterOptions{Exclude: []string{"B"}})
	if _, ok := got["B"]; ok {
		t.Error("B should be excluded")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestFilter_NoOptions(t *testing.T) {
	src := map[string]string{"X": "1", "Y": "2"}
	got := Filter(src, FilterOptions{})
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestFilter_PrefixAndExclude(t *testing.T) {
	src := map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"OTHER": "3",
	}
	got := Filter(src, FilterOptions{Prefix: "APP_", Exclude: []string{"APP_B"}})
	if len(got) != 1 || got["APP_A"] != "1" {
		t.Errorf("unexpected result: %v", got)
	}
}
