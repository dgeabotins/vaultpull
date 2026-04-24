package env

import (
	"strings"
	"testing"
)

func TestFlatten_NoChanges(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	res := Flatten(src, FlattenOptions{})
	if res.Renamed != 0 {
		t.Errorf("expected 0 renames, got %d", res.Renamed)
	}
	if res.Flattened["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST")
	}
}

func TestFlatten_ReplacesDots(t *testing.T) {
	src := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	res := Flatten(src, FlattenOptions{Separator: "_"})
	if res.Renamed != 2 {
		t.Errorf("expected 2 renames, got %d", res.Renamed)
	}
	if _, ok := res.Flattened["db_host"]; !ok {
		t.Error("expected key db_host to exist")
	}
	if _, ok := res.Flattened["db_port"]; !ok {
		t.Error("expected key db_port to exist")
	}
}

func TestFlatten_ReplacesDashes(t *testing.T) {
	src := map[string]string{
		"api-key": "secret",
	}
	res := Flatten(src, FlattenOptions{Separator: "_"})
	if _, ok := res.Flattened["api_key"]; !ok {
		t.Error("expected key api_key")
	}
	if res.Renamed != 1 {
		t.Errorf("expected 1 rename, got %d", res.Renamed)
	}
}

func TestFlatten_Uppercase(t *testing.T) {
	src := map[string]string{
		"db.host": "localhost",
	}
	res := Flatten(src, FlattenOptions{Uppercase: true})
	if _, ok := res.Flattened["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST after uppercase+flatten")
	}
}

func TestFlatten_Prefix(t *testing.T) {
	src := map[string]string{
		"HOST": "localhost",
	}
	res := Flatten(src, FlattenOptions{Prefix: "APP_"})
	if _, ok := res.Flattened["APP_HOST"]; !ok {
		t.Error("expected key APP_HOST")
	}
	if res.Renamed != 1 {
		t.Errorf("expected 1 rename, got %d", res.Renamed)
	}
}

func TestFlattenResult_Summary(t *testing.T) {
	src := map[string]string{"a.b": "1", "c": "2"}
	res := Flatten(src, FlattenOptions{})
	s := res.Summary()
	if !strings.Contains(s, "2 keys") {
		t.Errorf("summary missing key count: %s", s)
	}
}

func TestFlattenResult_Summary_NoRenames(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	res := Flatten(src, FlattenOptions{})
	s := res.Summary()
	if !strings.Contains(s, "no renames") {
		t.Errorf("expected 'no renames' in summary: %s", s)
	}
}
