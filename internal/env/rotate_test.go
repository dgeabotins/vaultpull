package env

import (
	"strings"
	"testing"
)

func TestRotate_AllKeys(t *testing.T) {
	current := map[string]string{"DB_PASS": "old", "API_KEY": "oldkey"}
	incoming := map[string]string{"DB_PASS": "new", "API_KEY": "newkey"}

	res, err := Rotate(current, incoming, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Fatalf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if current["DB_PASS"] != "new" {
		t.Errorf("expected DB_PASS=new, got %s", current["DB_PASS"])
	}
	if current["DB_PASS_PREVIOUS"] != "old" {
		t.Errorf("expected DB_PASS_PREVIOUS=old, got %s", current["DB_PASS_PREVIOUS"])
	}
	if current["API_KEY_PREVIOUS"] != "oldkey" {
		t.Errorf("expected API_KEY_PREVIOUS=oldkey, got %s", current["API_KEY_PREVIOUS"])
	}
}

func TestRotate_SelectedKeys(t *testing.T) {
	current := map[string]string{"DB_PASS": "old", "API_KEY": "oldkey"}
	incoming := map[string]string{"DB_PASS": "new", "API_KEY": "newkey"}

	res, err := Rotate(current, incoming, RotateOptions{Keys: []string{"DB_PASS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "DB_PASS" {
		t.Errorf("expected only DB_PASS rotated")
	}
	// API_KEY should be untouched
	if current["API_KEY"] != "oldkey" {
		t.Errorf("API_KEY should not have changed")
	}
}

func TestRotate_SkipsMissingIncoming(t *testing.T) {
	current := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"OTHER": "val"}

	res, err := Rotate(current, incoming, RotateOptions{Keys: []string{"DB_PASS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestRotate_CustomSuffix(t *testing.T) {
	current := map[string]string{"TOKEN": "abc"}
	incoming := map[string]string{"TOKEN": "xyz"}

	_, err := Rotate(current, incoming, RotateOptions{Suffix: "_OLD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current["TOKEN_OLD"] != "abc" {
		t.Errorf("expected TOKEN_OLD=abc, got %s", current["TOKEN_OLD"])
	}
}

func TestRotate_DryRun(t *testing.T) {
	current := map[string]string{"SECRET": "original"}
	incoming := map[string]string{"SECRET": "rotated"}

	res, err := Rotate(current, incoming, RotateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current["SECRET"] != "original" {
		t.Errorf("dry-run should not modify map")
	}
	if len(res.Rotated) != 1 {
		t.Errorf("dry-run should still report rotated keys")
	}
	if !strings.Contains(res.Summary(), "dry-run") {
		t.Errorf("summary should mention dry-run")
	}
}

func TestRotateSummary_NoChanges(t *testing.T) {
	res := RotateResult{Rotated: nil, Skipped: nil}
	if !strings.Contains(res.Summary(), "0") {
		t.Errorf("expected zero counts in summary")
	}
}
