package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRecoverMiddleware ensures panics in handlers are recovered and return HTTP 500
func TestRecoverMiddleware(t *testing.T) {
	h := recoverMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

// TestNewGraphQLHandlerBasic ensures the GraphQL handler can be created and responds to requests.
// This uses a nil repository to avoid needing a real database; resolvers should handle nil repo gracefully.
func TestNewGraphQLHandlerBasic(t *testing.T) {
	h := NewGraphQLHandler(nil)

	// Query the playground root (GraphQL handler expects POST for /query, but handler should respond to GET for introspection or similar)
	req := httptest.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Fatalf("handler returned 500 internal server error")
	}

	// Try a minimal POST (GraphQL empty body)
	postReq := httptest.NewRequest("POST", "/query", strings.NewReader("{}"))
	postReq.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	h.ServeHTTP(w2, postReq)

	resp2 := w2.Result()
	if resp2.StatusCode == http.StatusInternalServerError {
		t.Fatalf("handler POST returned 500 internal server error")
	}
}
