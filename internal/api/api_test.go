package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

func TestHealthIsPublicEvenOffLAN(t *testing.T) {
	e := newEnv(t)
	rec := e.do(req{method: "GET", path: "/api/health", remote: "10.9.9.9:1000"})
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
	e := newEnv(t)
	for path, op := range operations(e.api) {
		if op.OperationID == "" || op.Summary == "" {
			t.Errorf("%s lacks operationId or summary", path)
		}
	}
}

func TestOperationIDsAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for path, op := range operations(newEnv(t).api) {
		if seen[op.OperationID] {
			t.Errorf("duplicate operationId %s at %s", op.OperationID, path)
		}
		seen[op.OperationID] = true
	}
}

func TestCaseVariantPathIs404(t *testing.T) {
	e := newEnv(t)
	for _, p := range []string{"/API/health", "/Api/v1/auth/session", "/api/V1/setup"} {
		rec := e.do(req{method: "GET", path: p})
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", p, rec.Code)
		}
	}
}

func TestResponsesCarryNoSchemaLink(t *testing.T) {
	rec := newEnv(t).do(req{method: "GET", path: "/api/health"})
	if strings.Contains(rec.Body.String(), "$schema") {
		t.Fatalf("body links to a nonexistent schema route: %s", rec.Body.String())
	}
}

func TestNonLANSourceGets403ExceptHealth(t *testing.T) {
	e := newEnv(t)
	for _, remote := range []string{"203.0.113.9:5000", "192.168.2.20:5000", "[2001:db8::1]:5000", "garbage"} {
		for _, p := range []string{"/api/v1/setup", "/api/v1/auth/session", "/api/openapi.json", "/nope"} {
			rec := e.do(req{method: "GET", path: p, remote: remote})
			if rec.Code != http.StatusForbidden || errCode(t, rec) != "not_lan" {
				t.Errorf("%s from %s: %d %s", p, remote, rec.Code, rec.Body)
			}
		}
	}
	if rec := e.do(req{method: "POST", path: "/api/v1/setup/password", remote: "203.0.113.9:5000", body: `{"password":"longenough1"}`}); rec.Code != http.StatusForbidden {
		t.Fatalf("setup from uplink: %d", rec.Code)
	}
}

func TestOpenAPIIsPublicOnLAN(t *testing.T) {
	rec := newEnv(t).do(req{method: "GET", path: "/api/openapi.json"})
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"openapi":"3.1`) {
		t.Fatalf("%d %.80s", rec.Code, rec.Body)
	}
}

// operations lists every registered operation keyed by "METHOD path".
func operations(a huma.API) map[string]*huma.Operation {
	out := map[string]*huma.Operation{}
	for path, item := range a.OpenAPI().Paths {
		for method, op := range map[string]*huma.Operation{
			"GET": item.Get, "PUT": item.Put, "POST": item.Post, "DELETE": item.Delete, "PATCH": item.Patch,
		} {
			if op != nil {
				out[method+" "+path] = op
			}
		}
	}
	return out
}

func TestPanicBecomes500AndProcessSurvives(t *testing.T) {
	e := newEnv(t)
	e.acc.panicOnShadow = true
	rec := e.do(req{method: "GET", path: "/api/v1/setup"})
	if rec.Code != http.StatusInternalServerError || errCode(t, rec) != "internal" {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
	e.acc.panicOnShadow = false
	if rec := e.do(req{method: "GET", path: "/api/v1/setup"}); rec.Code != http.StatusOK {
		t.Fatalf("next request: %d", rec.Code)
	}
}

func TestStrictBodiesAre400NamingTheField(t *testing.T) {
	cases := []struct{ body, mention string }{
		{`{"password":"x","extra":1}`, "extra"},
		{`{"password":null}`, "password"},
		{`{"password":"x"} {"trailing":1}`, ""},
		{`{"Password":"x"}`, "Password"}, // keys are case-sensitive
		{`null`, ""},
	}
	for _, c := range cases {
		e := newEnv(t)
		rec := e.do(req{method: "POST", path: "/api/v1/auth/login", body: c.body})
		if rec.Code != http.StatusBadRequest || errCode(t, rec) != "invalid_request" {
			t.Errorf("%s: %d %s", c.body, rec.Code, rec.Body)
		}
		if c.mention != "" && !strings.Contains(rec.Body.String(), c.mention) {
			t.Errorf("%s: error does not name %q: %s", c.body, c.mention, rec.Body)
		}
		if strings.Contains(rec.Body.String(), `"title"`) || strings.Contains(rec.Body.String(), `"x"`) {
			t.Errorf("%s: huma shape or input value leaked: %s", c.body, rec.Body)
		}
	}
}

func TestUnknownRouteIsJSON404(t *testing.T) {
	rec := newEnv(t).do(req{method: "GET", path: "/api/v1/nothing"})
	if rec.Code != http.StatusNotFound || errCode(t, rec) != "not_found" {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
}
