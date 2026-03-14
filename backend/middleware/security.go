package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"
)

// allowedOriginsCache caches the parsed ALLOWED_ORIGINS value so we avoid
// strings.Split and slice allocation on every request. The value is computed
// lazily on first use and re-computed whenever the env var changes.
var (
	originsMu        sync.Mutex
	cachedOriginsEnv string
	cachedOrigins    []string
)

func getAllowedOrigins() []string {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")

	originsMu.Lock()
	defer originsMu.Unlock()

	if originsEnv == cachedOriginsEnv && cachedOrigins != nil {
		return cachedOrigins
	}
	cachedOriginsEnv = originsEnv
	if originsEnv == "" {
		cachedOrigins = []string{"http://localhost:8081", "http://localhost:19000"}
		return cachedOrigins
	}
	origins := strings.Split(originsEnv, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	cachedOrigins = result
	return cachedOrigins
}

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (only if TLS is enabled)
		if os.Getenv("ENABLE_TLS") == "true" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content Security Policy (restrictive for API)
		// Skip strict CSP in development to allow GraphQL playground to function
		if os.Getenv("ENV") != "development" {
			w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		}

		// Referrer policy
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Permissions policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

// CORS handles Cross-Origin Resource Sharing with configurable allowed origins
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := getAllowedOrigins()

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
