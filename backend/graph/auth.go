package graph

import (
	"context"
	"errors"

	"github.com/grillinr/nq/graph/model"
)

type contextKey string

const currentUserKey contextKey = "currentUser"

// WithCurrentUser sets the current user in context
func WithCurrentUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

// CurrentUser returns the current user from context
func CurrentUser(ctx context.Context) (*model.User, error) {
	value := ctx.Value(currentUserKey)
	if value == nil {
		return nil, errors.New("unauthorized")
	}

	user, ok := value.(*model.User)
	if !ok || user == nil {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
