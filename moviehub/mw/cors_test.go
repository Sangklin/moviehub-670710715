package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	for _, origin := range []string{"http://localhost:3000", "http://localhost:3001"} {
		t.Run(origin, func(t *testing.T) {
			called := false
			handler := CORS([]string{"http://localhost:3000"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			}))
			req := httptest.NewRequest(http.MethodOptions, "/api/reviews/1", nil)
			req.Header.Set("Origin", origin)
			req.Header.Set("Access-Control-Request-Method", "DELETE")
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != http.StatusNoContent || called {
				t.Fatalf("status=%d handler called=%v", res.Code, called)
			}
			if origin == "http://localhost:3000" {
				for key, value := range map[string]string{
					"Access-Control-Allow-Origin":  origin,
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Authorization, Content-Type",
					"Vary":                         "Origin",
				} {
					if got := res.Header().Get(key); got != value {
						t.Errorf("%s=%q want %q", key, got, value)
					}
				}
			} else if res.Header().Get("Access-Control-Allow-Origin") != "" || res.Header().Get("Access-Control-Allow-Headers") != "" {
				t.Fatal("unapproved origin received CORS permission")
			}
		})
	}
}
