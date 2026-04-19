package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForPatch(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPatch_SetNewKey(t *testing.T) {
	p := writeTempEnvForPatch(t, "FOO=bar\n")
	res, err := Patch(p, []PatchOp{{Key: "BAZ", Value: "qux"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Set) != 1 || res.Set[0] != "BAZ" {
		t.Errorf("expected BAZ in Set, got %v", res.Set)
	}
	m, _ := ToMap(p)
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", m["BAZ"])
	}
}

func TestPatch_UpdateExistingKey(t *testing.T) {
	p := writeTempEnvForPatch(t, "FOO=old\n")
	_, err := Patch(p, []PatchOp{{Key: "FOO", Value: "new"}})
	if err != nil {
		t.Fatal(err)
	}
	m, _ := ToMap(p)
	if m["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", m["FOO"])
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	p := writeTempEnvForPatch(t, "FOO=bar\nREMOVE=me\n")
	res, err := Patch(p, []PatchOp{{Key: "REMOVE", Delete: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Deleted) != 1 || res.Deleted[0] != "REMOVE" {
		t.Errorf("expected REMOVE in Deleted, got %v", res.Deleted)
	}
	m, _ := ToMap(p)
	if _, ok := m["REMOVE"]; ok {
		t.Error("expected REMOVE to be deleted")
	}
}

func TestPatch_NonExistentFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	_, err := Patch(p, []PatchOp{{Key: "NEW", Value: "val"}})
	if err != nil {
		t.Fatal(err)
	}
	m, _ := ToMap(p)
	if m["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %s", m["NEW"])
	}
}

func TestPatchResult_Summary(t *testing.T) {
	r := PatchResult{Set: []string{"A", "B"}, Deleted: []string{"C"}}
	s := r.Summary()
	if s != "2 set, 1 deleted" {
		t.Errorf("unexpected summary: %s", s)
	}
	empty := PatchResult{}
	if empty.Summary() != "no changes" {
		t.Errorf("expected 'no changes', got %s", empty.Summary())
	}
}
