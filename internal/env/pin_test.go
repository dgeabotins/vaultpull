package env

import (
	"testing"
)

func TestPin_AllKeys(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	out, res, err := Pin(env, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d", len(res.Pinned))
	}
	if !IsPinned(out["FOO"]) {
		t.Errorf("expected FOO to be pinned, got %q", out["FOO"])
	}
}

func TestPin_SelectedKeys(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	out, res, err := Pin(env, PinOptions{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 1 || res.Pinned[0] != "FOO" {
		t.Errorf("expected only FOO pinned, got %v", res.Pinned)
	}
	if IsPinned(out["BAZ"]) {
		t.Errorf("BAZ should not be pinned")
	}
}

func TestPin_SkipsAlreadyPinned(t *testing.T) {
	env := map[string]string{
		"FOO": "bar # pinned",
	}
	_, res, err := Pin(env, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Pinned) != 0 {
		t.Errorf("expected 0 pinned, got %d", len(res.Pinned))
	}
}

func TestPin_DryRun(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	out, res, err := Pin(env, PinOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 1 {
		t.Errorf("expected 1 in pinned report")
	}
	if IsPinned(out["KEY"]) {
		t.Errorf("dry-run should not modify value")
	}
}

func TestPin_PrefixFilter(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	_, res, err := Pin(env, PinOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d: %v", len(res.Pinned), res.Pinned)
	}
}

func TestIsPinned_AndUnpin(t *testing.T) {
	v := "myvalue # pinned"
	if !IsPinned(v) {
		t.Errorf("expected IsPinned to return true")
	}
	if got := UnpinValue(v); got != "myvalue" {
		t.Errorf("expected 'myvalue', got %q", got)
	}
}

func TestPin_NilMap(t *testing.T) {
	_, _, err := Pin(nil, PinOptions{})
	if err == nil {
		t.Error("expected error for nil map")
	}
}
