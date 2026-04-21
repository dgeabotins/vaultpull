package env

import (
	"testing"
)

func TestMergeMap_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "old", "B": "keep"}
	src := map[string]string{"A": "new", "C": "added"}

	res, err := MergeMap(dst, src, MergeOptions{Strategy: StrategyOverwrite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["A"] != "new" {
		t.Errorf("expected A=new, got %s", dst["A"])
	}
	if dst["C"] != "added" {
		t.Errorf("expected C=added, got %s", dst["C"])
	}
	if res.Updated != 1 {
		t.Errorf("expected 1 updated, got %d", res.Updated)
	}
	if res.Added != 1 {
		t.Errorf("expected 1 added, got %d", res.Added)
	}
}

func TestMergeMap_KeepExisting(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new", "B": "added"}

	res, err := MergeMap(dst, src, MergeOptions{Strategy: StrategyKeepExisting})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["A"] != "old" {
		t.Errorf("expected A=old (preserved), got %s", dst["A"])
	}
	if dst["B"] != "added" {
		t.Errorf("expected B=added, got %s", dst["B"])
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestMergeMap_ErrorOnConflict(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new"}

	res, err := MergeMap(dst, src, MergeOptions{Strategy: StrategyError})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "A" {
		t.Errorf("expected conflict on A, got %v", res.Conflicts)
	}
}

func TestMergeMap_PrefixFilter(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}

	res, err := MergeMap(dst, src, MergeOptions{Strategy: StrategyOverwrite, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("expected 2 added, got %d", res.Added)
	}
	if _, ok := dst["OTHER"]; ok {
		t.Error("OTHER should not have been merged")
	}
}

func TestMergeMap_IdenticalValuesSkipped(t *testing.T) {
	dst := map[string]string{"X": "same"}
	src := map[string]string{"X": "same"}

	res, err := MergeMap(dst, src, MergeOptions{Strategy: StrategyOverwrite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped for identical value, got %d", res.Skipped)
	}
	if res.Updated != 0 {
		t.Errorf("expected 0 updated, got %d", res.Updated)
	}
}
