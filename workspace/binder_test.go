package render

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func (t *testStruct) Bind(r *http.Request) error {
	if t.Field1 == "invalid" {
		return errors.New("validation failed")
	}
	return nil
}

func TestBindDecodingFailure(t *testing.T) {
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"field1": "value", "field2": "not-an-int"}`))
	r.Header.Set("Content-Type", "application/json")

	v := &testStruct{Field1: "default1", Field2: 42}
	err := Bind(r, v)
	if err == nil {
		t.Error("expected error during decoding, got nil")
	}

	if v.Field1 != "default1" || v.Field2 != 42 {
		t.Errorf("expected struct to remain unchanged, got %+v", v)
	}
}

func TestBindValidationFailure(t *testing.T) {
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"field1": "invalid", "field2": 100}`))
	r.Header.Set("Content-Type", "application/json")

	v := &testStruct{Field1: "default1", Field2: 42}
	err := Bind(r, v)
	if err == nil || err.Error() != "validation failed" {
		t.Errorf("expected validation error, got %v", err)
	}

	if v.Field1 != "default1" || v.Field2 != 42 {
		t.Errorf("expected struct to remain unchanged, got %+v", v)
	}
}

func TestBindSuccess(t *testing.T) {
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"field1": "valid", "field2": 100}`))
	r.Header.Set("Content-Type", "application/json")

	v := &testStruct{Field1: "default1", Field2: 42}
	err := Bind(r, v)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v.Field1 != "valid" || v.Field2 != 100 {
		t.Errorf("expected struct to be updated, got %+v", v)
	}
}
