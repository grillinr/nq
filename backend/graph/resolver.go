package graph

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/grillinr/nq/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const (
	// Maximum number of search statuses to keep in memory
	maxSearchStatuses = 10000
	// Time to keep completed search statuses before cleanup
	searchStatusTTL = 1 * time.Hour
	// Interval for periodic cleanup of old entries
	cleanupInterval = 15 * time.Minute
)

type Resolver struct {
	Repo           db.Repository
	searchStatuses *lru.Cache[string, searchStatus]
	searchMu       sync.RWMutex
	stopCleanup    chan struct{}
	cleanupDone    chan struct{}
}

// NewResolver creates a new resolver with database repository
func NewResolver(repo db.Repository) *Resolver {
	cache, err := lru.New[string, searchStatus](maxSearchStatuses)
	if err != nil {
		panic(err) // Should never happen with valid size
	}
	
	resolver := &Resolver{
		Repo:           repo,
		searchStatuses: cache,
		stopCleanup:    make(chan struct{}),
		cleanupDone:    make(chan struct{}),
	}
	
	// Start periodic cleanup goroutine
	go resolver.periodicCleanup()
	
	return resolver
}

// Close stops the periodic cleanup goroutine
func (r *Resolver) Close() {
	close(r.stopCleanup)
	<-r.cleanupDone
}

// periodicCleanup removes old completed search statuses
func (r *Resolver) periodicCleanup() {
	defer close(r.cleanupDone)
	
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			r.cleanupOldStatuses()
		case <-r.stopCleanup:
			return
		}
	}
}

// cleanupOldStatuses removes completed search statuses older than TTL
func (r *Resolver) cleanupOldStatuses() {
	r.searchMu.Lock()
	defer r.searchMu.Unlock()
	
	now := time.Now()
	var keysToRemove []string
	
	// Collect keys to remove
	for _, key := range r.searchStatuses.Keys() {
		if value, ok := r.searchStatuses.Peek(key); ok {
			if value.completedAt != nil {
				age := now.Sub(*value.completedAt)
				if age > searchStatusTTL {
					keysToRemove = append(keysToRemove, key)
				}
			}
		}
	}
	
	// Remove old entries
	for _, key := range keysToRemove {
		r.searchStatuses.Remove(key)
	}
}
