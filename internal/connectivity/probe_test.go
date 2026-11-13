package connectivity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPProbe(t *testing.T) {
	for _, tc := range []struct {
		code     int
		location string
		want     Internet
	}{
		{204, "", Online}, {302, "/hotel", Captive}, {301, "https://example.org/", Captive},
		{302, "ftp://example.org/", Offline}, {302, "http://[bad", Offline}, {302, "", Offline},
		{200, "", Offline}, {500, "", Offline},
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/check" || r.Header.Get("Cookie") != "" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(500)
				return
			}
			if tc.location != "" {
				w.Header().Set("Location", tc.location)
			}
			w.WriteHeader(tc.code)
		}))
		got := NewHTTPProbe(server.URL+"/check", time.Second)(context.Background())
		server.Close()
		if got != tc.want {
			t.Errorf("%d %q: got %s want %s", tc.code, tc.location, got, tc.want)
		}
	}
}

func TestHTTPProbeTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	if got := NewHTTPProbe(server.URL, 25*time.Millisecond)(context.Background()); got != Offline {
		t.Fatal(got)
	}
}

func TestHTTPProbeUsesFreshConnection(t *testing.T) {
	addresses := make(chan string, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addresses <- r.RemoteAddr
		w.WriteHeader(204)
	}))
	defer server.Close()
	probe := NewHTTPProbe(server.URL, time.Second)
	if probe(t.Context()) != Online || probe(t.Context()) != Online {
		t.Fatal("expected online")
	}
	if first, second := receive(t, addresses), receive(t, addresses); first == second {
		t.Fatal("probe reused a connection")
	}
}
