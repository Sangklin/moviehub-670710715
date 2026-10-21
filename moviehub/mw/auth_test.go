package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"moviehub/auth"
)

func TestLoginGate(t *testing.T) {
	issuer := &auth.Issuer{Secret: []byte("test-only-secret-that-is-long-enough"), TTL: time.Hour}
	token, err := issuer.Issue(7, "member")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, header string
		status       int
	}{
		{"no ticket", "", http.StatusUnauthorized},
		{"missing Bearer", token, http.StatusUnauthorized},
		{"altered ticket", "Bearer " + token + "x", http.StatusUnauthorized},
		{"valid ticket", "Bearer " + token, http.StatusNoContent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			handler := RequireLogin(issuer, func(w http.ResponseWriter, r *http.Request) {
				called = true
				if who := WhoFrom(r.Context()); who.ID != 7 || who.Role != "member" {
					t.Errorf("unexpected identity: %+v", who)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
			req.Header.Set("Authorization", tc.header)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != tc.status || called != (tc.status == http.StatusNoContent) {
				t.Fatalf("status=%d called=%v", res.Code, called)
			}
		})
	}
}

func TestAdminGate(t *testing.T) {
	issuer := &auth.Issuer{Secret: []byte("test-only-secret-that-is-long-enough"), TTL: time.Hour}
	for _, role := range []string{"member", "admin"} {
		t.Run(role, func(t *testing.T) {
			token, err := issuer.Issue(1, role)
			if err != nil {
				t.Fatal(err)
			}
			handler := RequireLogin(issuer, RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			req := httptest.NewRequest(http.MethodGet, "/api/admin/members", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			want := http.StatusForbidden
			if role == "admin" {
				want = http.StatusNoContent
			}
			if res.Code != want {
				t.Fatalf("got %d want %d", res.Code, want)
			}
		})
	}
}
