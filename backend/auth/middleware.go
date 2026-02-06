package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/grillinr/nq/graph"
	"github.com/grillinr/nq/graph/model"
)

type contextKey string

const contextKeyAuthClaims contextKey = "authClaims"

// AuthMiddleware validates access tokens and sets auth claims in context.
// It is permissive for unauthenticated requests and only blocks invalid tokens.
func AuthMiddleware(validator *Validator, repo ResolverRepo) func(http.Handler) http.Handler {
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
			if repo != nil {
				subject := getString(authUser.AuthSubject)
				if subject == "" {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte("missing subject"))
					return
				}
				user, err := repo.GetOrCreateUserByAuth(ctx, "auth0", subject, authUser.Email, authUser.Name, authUser.AvatarURL)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte("user lookup failed"))
					return
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
