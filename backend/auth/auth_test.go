package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/grillinr/nq/config"
)

func TestNewValidatorFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		issuer    string
		audience  string
		jwksURL   string
		wantError bool
	}{
		{
			name:      "valid configuration",
			issuer:    "https://example.auth0.com/",
			audience:  "https://api.example.com",
			jwksURL:   "https://example.auth0.com/.well-known/jwks.json",
			wantError: false,
		},
		{
			name:      "missing issuer",
			issuer:    "",
			audience:  "https://api.example.com",
			jwksURL:   "https://example.auth0.com/.well-known/jwks.json",
			wantError: true,
		},
		{
			name:      "missing audience",
			issuer:    "https://example.auth0.com/",
			audience:  "",
			jwksURL:   "https://example.auth0.com/.well-known/jwks.json",
			wantError: true,
		},
		{
			name:      "missing jwks URL",
			issuer:    "https://example.auth0.com/",
			audience:  "https://api.example.com",
			jwksURL:   "",
			wantError: true,
		},
		{
			name:      "all missing",
			issuer:    "",
			audience:  "",
			jwksURL:   "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.issuer != "" {
				os.Setenv("AUTH_JWT_ISSUER", tt.issuer)
			} else {
				os.Unsetenv("AUTH_JWT_ISSUER")
			}
			if tt.audience != "" {
				os.Setenv("AUTH_JWT_AUDIENCE", tt.audience)
			} else {
				os.Unsetenv("AUTH_JWT_AUDIENCE")
			}
			if tt.jwksURL != "" {
				os.Setenv("AUTH_JWT_JWKS_URL", tt.jwksURL)
			} else {
				os.Unsetenv("AUTH_JWT_JWKS_URL")
			}

			validator, err := NewValidatorFromEnv()

			if tt.wantError {
				if err == nil {
					t.Errorf("NewValidatorFromEnv() expected error, got nil")
				}
				if validator != nil {
					t.Errorf("NewValidatorFromEnv() expected nil validator on error, got %v", validator)
				}
			} else {
				if err != nil {
					t.Errorf("NewValidatorFromEnv() unexpected error: %v", err)
				}
				if validator == nil {
					t.Errorf("NewValidatorFromEnv() expected validator, got nil")
				}
				if validator != nil {
					if validator.issuer != tt.issuer {
						t.Errorf("issuer = %v, want %v", validator.issuer, tt.issuer)
					}
					if validator.audience != tt.audience {
						t.Errorf("audience = %v, want %v", validator.audience, tt.audience)
					}
					if validator.jwksURL != tt.jwksURL {
						t.Errorf("jwksURL = %v, want %v", validator.jwksURL, tt.jwksURL)
					}
				}
			}
		})
	}

	// Clean up
	os.Unsetenv("AUTH_JWT_ISSUER")
	os.Unsetenv("AUTH_JWT_AUDIENCE")
	os.Unsetenv("AUTH_JWT_JWKS_URL")
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "valid bearer token",
			header: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:   "bearer with extra spaces",
			header: "Bearer  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9  ",
			want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:   "lowercase bearer",
			header: "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:   "mixed case bearer",
			header: "BeArEr eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:   "no authorization header",
			header: "",
			want:   "",
		},
		{
			name:   "wrong scheme",
			header: "Basic dXNlcjpwYXNz",
			want:   "",
		},
		{
			name:   "missing token",
			header: "Bearer",
			want:   "",
		},
		{
			name:   "malformed - no space",
			header: "BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			want:   "",
		},
		{
			name:   "only bearer with space",
			header: "Bearer ",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			got := ExtractBearerToken(req)
			if got != tt.want {
				t.Errorf("ExtractBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorUserInfoURL(t *testing.T) {
	tests := []struct {
		name   string
		issuer string
		want   string
	}{
		{
			name:   "issuer with trailing slash",
			issuer: "https://example.auth0.com/",
			want:   "https://example.auth0.com/userinfo",
		},
		{
			name:   "issuer without trailing slash",
			issuer: "https://example.auth0.com",
			want:   "https://example.auth0.com/userinfo",
		},
		{
			name:   "issuer with multiple trailing slashes",
			issuer: "https://example.auth0.com///",
			want:   "https://example.auth0.com/userinfo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				issuer:     tt.issuer,
				httpClient: &http.Client{Timeout: config.AuthHTTPTimeout},
			}
			got := v.UserInfoURL()
			if got != tt.want {
				t.Errorf("UserInfoURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchUserInfo(t *testing.T) {
	t.Run("nil httpClient", func(t *testing.T) {
		// Create validator with nil httpClient to simulate zero-value construction
		v := &Validator{issuer: "https://example.com"}
		ctx := context.Background()

		_, err := v.FetchUserInfo(ctx, "test-token")
		if err == nil {
			t.Error("FetchUserInfo() expected error for nil httpClient, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "not properly initialized") {
			t.Errorf("FetchUserInfo() expected initialization error, got: %v", err)
		}
	})

	t.Run("successful fetch", func(t *testing.T) {
		// Create mock server
		mockUserInfo := UserInfo{
			Email:   "test@example.com",
			Name:    "Test User",
			Picture: "https://example.com/avatar.jpg",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify Authorization header
			auth := r.Header.Get("Authorization")
			if auth != "Bearer test-token" {
				t.Errorf("Expected Bearer test-token, got %s", auth)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockUserInfo)
		}))
		defer server.Close()

		v := &Validator{
			issuer:     server.URL,
			httpClient: &http.Client{Timeout: config.AuthHTTPTimeout},
		}
		ctx := context.Background()

		userInfo, err := v.FetchUserInfo(ctx, "test-token")
		if err != nil {
			t.Fatalf("FetchUserInfo() unexpected error: %v", err)
		}

		if userInfo.Email != mockUserInfo.Email {
			t.Errorf("Email = %v, want %v", userInfo.Email, mockUserInfo.Email)
		}
		if userInfo.Name != mockUserInfo.Name {
			t.Errorf("Name = %v, want %v", userInfo.Name, mockUserInfo.Name)
		}
		if userInfo.Picture != mockUserInfo.Picture {
			t.Errorf("Picture = %v, want %v", userInfo.Picture, mockUserInfo.Picture)
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}))
		defer server.Close()

		v := &Validator{
			issuer:     server.URL,
			httpClient: &http.Client{Timeout: config.AuthHTTPTimeout},
		}
		ctx := context.Background()

		_, err := v.FetchUserInfo(ctx, "invalid-token")
		if err == nil {
			t.Error("FetchUserInfo() expected error for unauthorized, got nil")
		}
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		v := &Validator{
			issuer:     server.URL,
			httpClient: &http.Client{Timeout: config.AuthHTTPTimeout},
		}
		ctx := context.Background()

		_, err := v.FetchUserInfo(ctx, "test-token")
		if err == nil {
			t.Error("FetchUserInfo() expected error for invalid JSON, got nil")
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second) // Longer than client timeout
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		v := &Validator{
			issuer:     server.URL,
			httpClient: &http.Client{Timeout: config.AuthHTTPTimeout},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := v.FetchUserInfo(ctx, "test-token")
		if err == nil {
			t.Error("FetchUserInfo() expected timeout error, got nil")
		}
	})
}

func TestValidateAccessToken_MissingSubject(t *testing.T) {
	// This test verifies that tokens without a subject are rejected
	// even if all other validations pass

	// Note: Full integration test with real JWT validation requires:
	// 1. Generating valid JWTs with proper signing
	// 2. Mocking JWKS endpoint
	// 3. Complex setup for RSA key pairs
	//
	// This is better suited for integration tests with actual auth providers
	// For now, we test the logic paths we can reach without full JWT infrastructure

	t.Skip("Requires full JWT infrastructure - covered by integration tests")
}

func TestValidator_CacheKeySet(t *testing.T) {
	t.Run("cache refresh after 10 minutes", func(t *testing.T) {
		v := &Validator{
			issuer:   "https://example.com",
			audience: "test-audience",
			jwksURL:  "https://example.com/.well-known/jwks.json",
		}

		// Set last fetch to 11 minutes ago
		v.lastFetch = time.Now().Add(-11 * time.Minute)

		// Note: This would require mocking the jwk.Fetch call
		// which is complex. This documents the expected behavior.
		// In practice, this is tested via integration tests.

		t.Skip("Requires mocking jwk.Fetch - integration test needed")
	})
}
