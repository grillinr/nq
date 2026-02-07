package graph

import (
	"sync"

	"github.com/grillinr/nq/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo           db.Repository
	searchStatuses map[string]searchStatus
	searchMu       sync.RWMutex
}

// NewResolver creates a new resolver with database repository
func NewResolver(repo db.Repository) *Resolver {
	return &Resolver{
		Repo:           repo,
		searchStatuses: make(map[string]searchStatus),
	}
}
