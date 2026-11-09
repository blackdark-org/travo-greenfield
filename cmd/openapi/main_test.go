package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/BlackDark/travo-greenfield/internal/api"
)

// The committed spec must describe exactly what the daemon serves.
func TestOutputEqualsServedSpec(t *testing.T) {
	var out bytes.Buffer
	if err := render(&out); err != nil {
		t.Fatal(err)
	}
	h, _ := api.New(specVersion)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var got, want any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("printed spec differs from served spec\nprinted: %s\nserved: %s", out.String(), rec.Body.String())
	}
}
