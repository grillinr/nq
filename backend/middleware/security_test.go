package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with security headers middleware
	handler := SecurityHeaders(nextHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute request
	handler.ServeHTTP(rec, req)

	// Check headers
	headers := rec.Header()

	if got := headers.Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("X-Frame-Options = %q, want %q", got, "DENY")
	}

	if got := headers.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want %q", got, "nosniff")
	}

	if got := headers.Get("X-XSS-Protection"); got != "1; mode=block" {
		t.Errorf("X-XSS-Protection = %q, want %q", got, "1; mode=block")
	}

	if got := headers.Get("Referrer-Policy"); got != "no-referrer" {
		t.Errorf("Referrer-Policy = %q, want %q", got, "no-referrer")
	}
}

func TestSecurityHeadersWithTLS(t *testing.T) {
	// Set TLS enabled
	os.Setenv("ENABLE_TLS", "true")
	defer os.Unsetenv("ENABLE_TLS")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := SecurityHeaders(nextHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	hsts := rec.Header().Get("Strict-Transport-Security")
	if hsts == "" {
		t.Error("HSTS header should be set when TLS is enabled")
	}
	expected := "max-age=31536000; includeSubDomains"
	if hsts != expected {
		t.Errorf("HSTS = %q, want %q", hsts, expected)
	}
}

func TestCORS_AllowedOrigin(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,https://app.example.com")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(nextHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	allowedOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowedOrigin != "https://example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", allowedOrigin, "https://example.com")
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "https://example.com")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(nextHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should not set CORS header for disallowed origin
	allowedOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowedOrigin == "https://evil.com" {
		t.Error("Should not allow origin from evil.com")
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:8081")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called for OPTIONS request")
	})

	handler := CORS(nextHandler)
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Preflight should return 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST, OPTIONS" {
		t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, "GET, POST, OPTIONS")
	}
}
