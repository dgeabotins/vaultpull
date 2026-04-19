package env

import (
	"strings"
	"testing"
)

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "1"}
	res := Compare(a, b)
	if _, ok := res.OnlyInA["BAR"]; !ok {
		t.Error("expected BAR in OnlyInA")
	}
	if len(res.OnlyInB) != 0 {
		t.Error("expected OnlyInB empty")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1", "BAZ": "3"}
	res := Compare(a, b)
	if _, ok := res.OnlyInB["BAZ"]; !ok {
		t.Error("expected BAZ in OnlyInB")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	res := Compare(a, b)
	pair, ok := res.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected values: %v", pair)
	}
}

func TestCompare_Same(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "1"}
	res := Compare(a, b)
	if _, ok := res.Same["X"]; !ok {
		t.Error("expected X in Same")
	}
}

func TestCompare_Summary(t *testing.T) {
	a := map[string]string{"A": "1", "B": "old"}
	b := map[string]string{"B": "new", "C": "3"}
	res := Compare(a, b)
	s := res.Summary()
	if !strings.Contains(s, "only-in-a=1") || !strings.Contains(s, "only-in-b=1") || !strings.Contains(s, "changed=1") {
		t.Errorf("unexpected summary: %s", s)
	}
}
