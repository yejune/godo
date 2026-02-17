package hook

import "testing"

func Test_GetStringField_found(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"mode": "do",
	}
	got := GetStringField(data, "name", "default")
	if got != "test" {
		t.Errorf("got %q, want %q", got, "test")
	}
}

func Test_GetStringField_missing_key(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
	}
	got := GetStringField(data, "missing", "fallback")
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func Test_GetStringField_empty_string(t *testing.T) {
	data := map[string]interface{}{
		"name": "",
	}
	got := GetStringField(data, "name", "default")
	if got != "default" {
		t.Errorf("got %q, want %q (empty string should return fallback)", got, "default")
	}
}

func Test_GetStringField_non_string_value(t *testing.T) {
	data := map[string]interface{}{
		"count": 42,
	}
	got := GetStringField(data, "count", "default")
	if got != "default" {
		t.Errorf("got %q, want %q (non-string should return fallback)", got, "default")
	}
}

func Test_GetStringField_nil_map(t *testing.T) {
	got := GetStringField(nil, "key", "fallback")
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}
