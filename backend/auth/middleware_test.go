package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

// mockValidator is a test double for Validator
type mockValidator struct {
	validateFunc func(ctx context.Context, token string) (*Claims, error)
	userInfoFunc func(ctx context.Context, token string) (*UserInfo, error)
}

func (m *mockValidator) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, token)
	}
	return nil, errors.New("not implemented")
}

func (m *mockValidator) FetchUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	if m.userInfoFunc != nil {
		return m.userInfoFunc(ctx, token)
	}
	return nil, errors.New("not implemented")
}

// mockResolverRepo is a test double for ResolverRepo
type mockResolverRepo struct {
	getOrCreateFunc func(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error)
}

func (m *mockResolverRepo) GetOrCreateUserByAuth(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error) {
	if m.getOrCreateFunc != nil {
		return m.getOrCreateFunc(ctx, provider, subject, email, name, avatarURL)
	}
	return nil, errors.New("not implemented")
}

func TestAuthMiddleware_NoValidator(t *testing.T) {
	// When validator is nil, middleware should pass through
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := AuthMiddleware(nil, nil)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "success" {
		t.Errorf("expected 'success', got %s", rec.Body.String())
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	// When no token is provided, middleware should pass through
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	mockVal := &mockValidator{}
	middleware := AuthMiddleware(mockVal, nil)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// When token validation fails, should return 401
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for invalid token")
	})

	mockVal := &mockValidator{
		validateFunc: func(ctx context.Context, token string) (*Claims, error) {
			return nil, errors.New("invalid token")
		},
	}

	middleware := AuthMiddleware(mockVal, nil)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
	if !contains(rec.Body.String(), "invalid token") {
		t.Errorf("expected 'invalid token' message, got %s", rec.Body.String())
	}
}

func TestAuthMiddleware_ValidToken_NoRepo(t *testing.T) {
	// Valid token without repo should pass claims to context
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Verify claims are in context
		user, err := CurrentAuthUser(r.Context())
		if err != nil {
			t.Errorf("expected user in context, got error: %v", err)
		}
		if user.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", user.Email)
		}

		w.WriteHeader(http.StatusOK)
	})

	mockVal := &mockValidator{
		validateFunc: func(ctx context.Context, token string) (*Claims, error) {
			return &Claims{
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/avatar.jpg",
			}, nil
		},
	}

	middleware := AuthMiddleware(mockVal, nil)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidToken_WithRepo(t *testing.T) {
	// Valid token with repo should create/fetch user
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	mockVal := &mockValidator{
		validateFunc: func(ctx context.Context, token string) (*Claims, error) {
			claims := &Claims{
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/avatar.jpg",
			}
			claims.Subject = "auth0|123456"
			return claims, nil
		},
	}

	mockRepo := &mockResolverRepo{
		getOrCreateFunc: func(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error) {
			if provider != "auth0" {
				t.Errorf("expected provider auth0, got %s", provider)
			}
			if subject != "auth0|123456" {
				t.Errorf("expected subject auth0|123456, got %s", subject)
			}
			if email != "test@example.com" {
				t.Errorf("expected email test@example.com, got %s", email)
			}

			return &model.User{
				ID:    uuid.New(),
				Email: email,
				Name:  name,
			}, nil
		},
	}

	middleware := AuthMiddleware(mockVal, mockRepo)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_MissingSubject(t *testing.T) {
	// Token without subject should be rejected when repo is provided
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for missing subject")
	})

	mockVal := &mockValidator{
		validateFunc: func(ctx context.Context, token string) (*Claims, error) {
			return &Claims{
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/avatar.jpg",
				// Subject is missing
			}, nil
		},
	}

	mockRepo := &mockResolverRepo{}

	middleware := AuthMiddleware(mockVal, mockRepo)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token-without-subject")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
	if !contains(rec.Body.String(), "missing subject") {
		t.Errorf("expected 'missing subject' message, got %s", rec.Body.String())
	}
}

func TestAuthMiddleware_RepoError(t *testing.T) {
	// Repo error should result in 401
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when repo fails")
	})

	mockVal := &mockValidator{
		validateFunc: func(ctx context.Context, token string) (*Claims, error) {
			claims := &Claims{
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/avatar.jpg",
			}
			claims.Subject = "auth0|123456"
			return claims, nil
		},
	}

	mockRepo := &mockResolverRepo{
		getOrCreateFunc: func(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error) {
			return nil, errors.New("database error")
		},
	}

	middleware := AuthMiddleware(mockVal, mockRepo)
	wrapped := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
	if !contains(rec.Body.String(), "user lookup failed") {
		t.Errorf("expected 'user lookup failed' message, got %s", rec.Body.String())
	}
}

func TestCurrentAuthUser(t *testing.T) {
	t.Run("valid claims in context", func(t *testing.T) {
		claims := &Claims{
			Email:   "test@example.com",
			Name:    "Test User",
			Picture: "https://example.com/avatar.jpg",
		}
		claims.Subject = "auth0|123456"

		ctx := context.WithValue(context.Background(), contextKeyAuthClaims, claims)

		user, err := CurrentAuthUser(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user.Email != "test@example.com" {
			t.Errorf("email = %v, want test@example.com", user.Email)
		}
		if user.Name != "Test User" {
			t.Errorf("name = %v, want Test User", user.Name)
		}
		if user.AvatarURL == nil || *user.AvatarURL != "https://example.com/avatar.jpg" {
			t.Errorf("avatarURL = %v, want https://example.com/avatar.jpg", user.AvatarURL)
		}
		if user.AuthProvider == nil || *user.AuthProvider != "auth0" {
			t.Errorf("authProvider = %v, want auth0", user.AuthProvider)
		}
		if user.AuthSubject == nil || *user.AuthSubject != "auth0|123456" {
			t.Errorf("authSubject = %v, want auth0|123456", user.AuthSubject)
		}
	})

	t.Run("no claims in context", func(t *testing.T) {
		ctx := context.Background()

		_, err := CurrentAuthUser(ctx)
		if err == nil {
			t.Error("expected error for missing claims, got nil")
		}
	})

	t.Run("nil claims in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextKeyAuthClaims, (*Claims)(nil))

		_, err := CurrentAuthUser(ctx)
		if err == nil {
			t.Error("expected error for nil claims, got nil")
		}
	})

	t.Run("empty picture becomes nil", func(t *testing.T) {
		claims := &Claims{
			Email:   "test@example.com",
			Name:    "Test User",
			Picture: "   ", // Empty/whitespace
		}
		claims.Subject = "auth0|123456"

		ctx := context.WithValue(context.Background(), contextKeyAuthClaims, claims)

		user, err := CurrentAuthUser(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user.AvatarURL != nil {
			t.Errorf("expected nil avatarURL for empty picture, got %v", *user.AvatarURL)
		}
	})
}

func TestStringPointer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *string
	}{
		{
			name:  "non-empty string",
			input: "test",
			want:  stringPtr("test"),
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "   ",
			want:  nil,
		},
		{
			name:  "string with spaces",
			input: " test ",
			want:  stringPtr(" test "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringPointer(tt.input)
			if (got == nil) != (tt.want == nil) {
				t.Errorf("stringPointer() nil mismatch: got %v, want %v", got, tt.want)
				return
			}
			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("stringPointer() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "non-nil pointer",
			input: stringPtr("test"),
			want:  "test",
		},
		{
			name:  "nil pointer",
			input: nil,
			want:  "",
		},
		{
			name:  "empty string pointer",
			input: stringPtr(""),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getString(tt.input)
			if got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}

func stringPtr(s string) *string {
	return &s
}
