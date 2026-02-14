package middleware

import (
	"log"
	"net/http"
	"os"
)

// ErrorSanitizer wraps responses to sanitize error messages in production
type ErrorSanitizer struct {
	http.ResponseWriter
	statusCode int
}

func (es *ErrorSanitizer) WriteHeader(code int) {
	es.statusCode = code
	es.ResponseWriter.WriteHeader(code)
}

// SanitizeErrors wraps the response writer to intercept error responses
func SanitizeErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In development, pass through errors unchanged
		env := os.Getenv("ENV")
		if env == "" || env == "development" {
			next.ServeHTTP(w, r)
			return
		}

		// In production, wrap the response writer
		wrapper := &ErrorSanitizer{ResponseWriter: w, statusCode: http.StatusOK}

		// Recover from panics and return sanitized error
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic recovered: %v", rec)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(wrapper, r)
	})
}
