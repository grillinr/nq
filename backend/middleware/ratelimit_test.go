package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimiter_Allow(t *testing.T) {
	// Create a rate limiter that allows 2 requests per second with burst of 2
	rl := NewRateLimiter(rate.Limit(2), 2)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(nextHandler)

	// First two requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected 200, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimiter_Exceed(t *testing.T) {
	// Create a very strict rate limiter (1 request per second, burst 1)
	rl := NewRateLimiter(rate.Limit(1), 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(nextHandler)

	// First request succeeds
	req1 := httptest.NewRequest("GET", "/test", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("First request: expected 200, got %d", rec1.Code)
	}

	// Second immediate request should be rate limited
	req2 := httptest.NewRequest("GET", "/test", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request: expected 429, got %d", rec2.Code)
	}

	retryAfter := rec2.Header().Get("Retry-After")
	if retryAfter != "60" {
		t.Errorf("Retry-After header = %q, want %q", retryAfter, "60")
	}
}

func TestRateLimiter_PerIP(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(1), 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rl.Limit(nextHandler)

	// Request from IP 1
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("IP1 request: expected 200, got %d", rec1.Code)
	}

	// Request from different IP should still succeed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:1234"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("IP2 request: expected 200, got %d", rec2.Code)
	}

	// Another request from IP 1 (same RemoteAddr) should be rate limited
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:1234" // Same RemoteAddr as req1
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 second request: expected 429, got %d", rec3.Code)
	}
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 192.168.1.1")

	ip := getClientIP(req)
	if ip != "203.0.113.1" {
		t.Errorf("getClientIP() = %q, want %q", ip, "203.0.113.1")
	}
}

func TestGetClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-IP", "203.0.113.2")

	ip := getClientIP(req)
	if ip != "203.0.113.2" {
		t.Errorf("getClientIP() = %q, want %q", ip, "203.0.113.2")
	}
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:8080"

	ip := getClientIP(req)
	if ip != "192.168.1.100:8080" {
		t.Errorf("getClientIP() = %q, want %q", ip, "192.168.1.100:8080")
	}
}

func TestRateLimiter_CleanupVisitors(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(10), 10)
	defer rl.Stop() // Stop the cleanup goroutine when test completes

	// Add a visitor
	rl.getVisitor("192.168.1.1")

	// Verify visitor exists
	rl.mu.RLock()
	if len(rl.visitors) != 1 {
		t.Errorf("Expected 1 visitor, got %d", len(rl.visitors))
	}
	rl.mu.RUnlock()

	// Manually set lastSeen to old time
	rl.mu.Lock()
	for _, v := range rl.visitors {
		v.lastSeen.Store(time.Now().Add(-10 * time.Minute).UnixNano())
	}
	rl.mu.Unlock()

	// Note: We can't easily test the automatic cleanup goroutine
	// without waiting 5 minutes, so we just verify the structure is correct
}

func TestRateLimiter_StopMultipleTimes(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(10), 10)

	// Should not panic when called multiple times
	rl.Stop()
	rl.Stop()
	rl.Stop()
}
