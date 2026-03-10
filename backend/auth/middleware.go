package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/grillinr/nq/graph"
	"github.com/grillinr/nq/graph/model"
)

type contextKey string

const contextKeyAuthClaims contextKey = "authClaims"

// userCacheEntry holds a cached user and its expiry time.
type userCacheEntry struct {
	user      *model.User
	expiresAt time.Time
}

// userCache is a simple in-memory cache for resolved users keyed by Auth0 subject.
// This avoids a DB MERGE write on every authenticated request.
type userCache struct {
	mu      sync.RWMutex
	entries map[string]userCacheEntry
}

const userCacheTTL = 60 * time.Second

func (c *userCache) get(subject string) (*model.User, bool) {
	c.mu.RLock()
	entry, ok := c.entries[subject]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.user, true
}

func (c *userCache) set(subject string, user *model.User) {
	c.mu.Lock()
	c.entries[subject] = userCacheEntry{user: user, expiresAt: time.Now().Add(userCacheTTL)}
	c.mu.Unlock()
}

// AuthMiddleware validates access tokens and sets auth claims in context.
// It is permissive for unauthenticated requests and only blocks invalid tokens.
func AuthMiddleware(validator TokenValidator, repo ResolverRepo) func(http.Handler) http.Handler {
	cache := &userCache{entries: make(map[string]userCacheEntry)}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if validator == nil {
				next.ServeHTTP(w, r)
				return
			}

			token := ExtractBearerToken(r)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := validator.ValidateAccessToken(r.Context(), token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyAuthClaims, claims)
			authUser, err := CurrentAuthUser(ctx)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("invalid user"))
				return
			}
			if authUser.Email == "" || authUser.Name == "" || authUser.AvatarURL == nil {
				if info, err := validator.FetchUserInfo(ctx, token); err == nil && info != nil {
					if authUser.Email == "" {
						authUser.Email = info.Email
					}
					if authUser.Name == "" {
						authUser.Name = info.Name
					}
					if authUser.AvatarURL == nil && info.Picture != "" {
						authUser.AvatarURL = stringPointer(info.Picture)
					}
				}
			}
			if repo != nil {
				subject := getString(authUser.AuthSubject)
				if subject == "" {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte("missing subject"))
					return
				}
				if authUser.Email == "" {
					authUser.Email = subject
				}
				// Check the in-memory cache before hitting the database.
				user, cached := cache.get(subject)
				if !cached {
					user, err = repo.GetOrCreateUserByAuth(ctx, "auth0", subject, authUser.Email, authUser.Name, authUser.AvatarURL)
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						_, _ = w.Write([]byte("user lookup failed"))
						return
					}
					cache.set(subject, user)
				}
				ctx = graph.WithCurrentUser(ctx, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CurrentAuthUser builds a user model from validated token claims.
func CurrentAuthUser(ctx context.Context) (*model.User, error) {
	value := ctx.Value(contextKeyAuthClaims)
	claims, ok := value.(*Claims)
	if !ok || claims == nil {
		return nil, errors.New("unauthorized")
	}

	avatar := strings.TrimSpace(claims.Picture)
	var avatarURL *string
	if avatar != "" {
		avatarURL = &avatar
	}

	return &model.User{
		Name:         claims.Name,
		Email:        claims.Email,
		AuthProvider: stringPointer("auth0"),
		AuthSubject:  stringPointer(claims.Subject),
		AvatarURL:    avatarURL,
	}, nil
}

type ResolverRepo interface {
	GetOrCreateUserByAuth(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error)
}

func getString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringPointer(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}
