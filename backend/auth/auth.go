package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Claims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

type Validator struct {
	issuer    string
	audience  string
	jwksURL   string
	keySet    jwk.Set
	lastFetch time.Time
}

type UserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewValidatorFromEnv() (*Validator, error) {
	issuer := os.Getenv("AUTH_JWT_ISSUER")
	audience := os.Getenv("AUTH_JWT_AUDIENCE")
	jwksURL := os.Getenv("AUTH_JWT_JWKS_URL")

	if issuer == "" || audience == "" || jwksURL == "" {
		return nil, errors.New("missing auth env: AUTH_JWT_ISSUER, AUTH_JWT_AUDIENCE, AUTH_JWT_JWKS_URL")
	}

	return &Validator{issuer: issuer, audience: audience, jwksURL: jwksURL}, nil
}

func (v *Validator) fetchKeySet(ctx context.Context) error {
	if v.keySet != nil && time.Since(v.lastFetch) < 10*time.Minute {
		return nil
	}

	keySet, err := jwk.Fetch(ctx, v.jwksURL)
	if err != nil {
		return err
	}

	v.keySet = keySet
	v.lastFetch = time.Now()
	return nil
}

func (v *Validator) ValidateAccessToken(ctx context.Context, tokenString string) (*Claims, error) {
	if err := v.fetchKeySet(ctx); err != nil {
		return nil, err
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}

		key, ok := v.keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key not found for kid: %s", kid)
		}

		var rawKey any
		if err := key.Raw(&rawKey); err != nil {
			return nil, err
		}
		return rawKey, nil
	}

	claims := &Claims{}
	parser := jwt.NewParser(
		jwt.WithAudience(v.audience),
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{"RS256"}),
	)
	_, err := parser.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims.Subject == "" {
		return nil, errors.New("missing subject")
	}

	return claims, nil
}

func (v *Validator) UserInfoURL() string {
	return strings.TrimRight(v.issuer, "/") + "/userinfo"
}

func (v *Validator) FetchUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, v.UserInfoURL(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("userinfo failed: %s", strings.TrimSpace(string(body)))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
