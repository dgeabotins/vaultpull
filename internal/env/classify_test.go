package env

import (
	"testing"
)

func TestClassify_URL(t *testing.T) {
	m := map[string]string{"DATABASE_URL": "https://example.com/db"}
	res := Classify(m)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Category != CategoryURL {
		t.Errorf("expected url, got %s", res[0].Category)
	}
}

func TestClassify_Secret(t *testing.T) {
	m := map[string]string{"API_SECRET": "supersecretvalue"}
	res := Classify(m)
	if res[0].Category != CategorySecret {
		t.Errorf("expected secret, got %s", res[0].Category)
	}
}

func TestClassify_Boolean(t *testing.T) {
	for _, v := range []string{"true", "false", "yes", "no"} {
		m := map[string]string{"FLAG": v}
		res := Classify(m)
		if res[0].Category != CategoryBoolean {
			t.Errorf("value %q: expected boolean, got %s", v, res[0].Category)
		}
	}
}

func TestClassify_Integer(t *testing.T) {
	m := map[string]string{"PORT": "8080"}
	res := Classify(m)
	if res[0].Category != CategoryInteger {
		t.Errorf("expected integer, got %s", res[0].Category)
	}
}

func TestClassify_Float(t *testing.T) {
	m := map[string]string{"RATIO": "3.14"}
	res := Classify(m)
	if res[0].Category != CategoryFloat {
		t.Errorf("expected float, got %s", res[0].Category)
	}
}

func TestClassify_Path(t *testing.T) {
	m := map[string]string{"CONFIG_PATH": "/etc/app/config.yaml"}
	res := Classify(m)
	if res[0].Category != CategoryPath {
		t.Errorf("expected path, got %s", res[0].Category)
	}
}

func TestClassify_JSON(t *testing.T) {
	m := map[string]string{"SETTINGS": `{"key":"val"}`}
	res := Classify(m)
	if res[0].Category != CategoryJSON {
		t.Errorf("expected json, got %s", res[0].Category)
	}
}

func TestClassify_Empty(t *testing.T) {
	m := map[string]string{"EMPTY_VAR": ""}
	res := Classify(m)
	if res[0].Category != CategoryEmpty {
		t.Errorf("expected empty, got %s", res[0].Category)
	}
}

func TestClassify_Unknown(t *testing.T) {
	m := map[string]string{"APP_NAME": "myapp"}
	res := Classify(m)
	if res[0].Category != CategoryUnknown {
		t.Errorf("expected unknown, got %s", res[0].Category)
	}
}

func TestClassify_MultipleKeys_Sorted(t *testing.T) {
	m := map[string]string{
		"Z_VAR": "hello",
		"A_VAR": "123",
	}
	res := Classify(m)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if res[0].Key != "A_VAR" {
		t.Errorf("expected A_VAR first, got %s", res[0].Key)
	}
	if res[0].Category != CategoryInteger {
		t.Errorf("expected integer for A_VAR, got %s", res[0].Category)
	}
}
