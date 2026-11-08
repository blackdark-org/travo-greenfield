package api

import (
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	h, _ := New("1.2.3")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var body struct{ Status, Version string }
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" || body.Version != "1.2.3" {
		t.Fatalf("body %+v", body)
	}
}

func TestEveryOperationHasIDAndSummary(t *testing.T) {
	_, a := New("test")
	for path, item := range a.OpenAPI().Paths {
		ops := map[string]*huma.Operation{
			"GET": item.Get, "PUT": item.Put, "POST": item.Post,
			"DELETE": item.Delete, "PATCH": item.Patch,
		}
		for method, op := range ops {
			if op != nil && (op.OperationID == "" || op.Summary == "") {
				t.Errorf("%s %s lacks operationId or summary", method, path)
			}
		}
	}
}

func TestCaseVariantPathIs404(t *testing.T) {
	h, _ := New("test")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/API/health", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d, want 404", rec.Code)
	}
}
