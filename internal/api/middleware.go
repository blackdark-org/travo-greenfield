package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/netip"

	"github.com/BlackDark/travo-greenfield/internal/auth"
	"github.com/danielgtaylor/huma/v2"
)

const (
	cookieName = "tg_session"
	csrfHeader = "X-CSRF-Token"
	publicKey  = "public"
)

type (
	infoKey   struct{}
	authedKey struct{}
)

// reqInfo is what handlers need from the transport.
type reqInfo struct {
	ip  string // unmapped, no port; the rate-limit key
	tls bool
}

// authed is the validated session of this request.
type authed struct {
	token string
	sess  auth.Session
}

func infoOf(ctx context.Context) reqInfo { v, _ := ctx.Value(infoKey{}).(reqInfo); return v }

func authedOf(ctx context.Context) authed { v, _ := ctx.Value(authedKey{}).(authed); return v }

func remoteAddr(hostport string) (netip.Addr, bool) {
	ap, err := netip.ParseAddrPort(hostport)
	if err != nil {
		return netip.Addr{}, false
	}
	return ap.Addr().Unmap(), true
}

// outer is the http-level chain: recover -> trust -> mux. Trust sits here so that 404s
// and the spec are LAN-only too; only /api/health answers everyone.
func (s *server) outer(mux http.Handler) http.Handler {
	return recoverPanics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			ip, ok := remoteAddr(r.RemoteAddr)
			if !ok || s.trusted == nil || !s.trusted(ip) {
				writeJSONError(w, http.StatusForbidden, "not_lan", "admin is reachable from the LAN only")
				return
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		mux.ServeHTTP(&jsonErrors{ResponseWriter: w}, r)
	}))
}

func writeJSONError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(&apiError{Status: status, Message: msg, Code: code})
}

// recoverPanics keeps the process alive when a handler panics (invariants §1).
func recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				slog.Error("handler panic", "method", r.Method, "path", r.URL.Path, "panic", v)
				writeJSONError(w, http.StatusInternalServerError, "internal", "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func isUnsafe(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

func isPublic(op *huma.Operation) bool {
	pub, _ := op.Metadata[publicKey].(bool)
	return pub // absent = not public: default-deny
}

// gate is the Huma middleware: request info, then auth, then CSRF. Routes are protected
// unless registered as public.
func (s *server) gate(ctx huma.Context, next func(huma.Context)) {
	ip, _ := remoteAddr(ctx.RemoteAddr())
	ctx = huma.WithValue(ctx, infoKey{}, reqInfo{ip: ip.String(), tls: ctx.TLS() != nil})
	if isPublic(ctx.Operation()) {
		next(ctx)
		return
	}
	c, err := huma.ReadCookie(ctx, cookieName)
	if err != nil {
		writeError(ctx, http.StatusUnauthorized, "unauthenticated", "login required")
		return
	}
	sess, ok := s.sessions.Lookup(c.Value, s.clock.Now())
	if !ok {
		writeError(ctx, http.StatusUnauthorized, "unauthenticated", "login required")
		return
	}
	if isUnsafe(ctx.Method()) && subtle.ConstantTimeCompare([]byte(ctx.Header(csrfHeader)), []byte(sess.CSRF)) != 1 {
		writeError(ctx, http.StatusForbidden, "csrf", "missing or wrong "+csrfHeader)
		return
	}
	next(huma.WithValue(ctx, authedKey{}, authed{token: c.Value, sess: sess}))
}

// sessionCookie builds tg_session; Secure only when the request came over TLS (the
// daemon serves plain HTTP on the LAN today).
func sessionCookie(token string, info reqInfo, maxAge int) http.Cookie {
	return http.Cookie{
		Name: cookieName, Value: token, Path: "/", MaxAge: maxAge,
		HttpOnly: true, Secure: info.tls, SameSite: http.SameSiteStrictMode,
	}
}

func clearedCookie(info reqInfo) http.Cookie { return sessionCookie("", info, -1) }
