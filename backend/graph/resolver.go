package graph

import (
	"nq/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo db.Repository
}

// NewResolver creates a new resolver with database repository
func NewResolver(repo db.Repository) *Resolver {
	return &Resolver{
		Repo: repo,
	}
}
