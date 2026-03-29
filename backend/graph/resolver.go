package graph

import (
	"fmt"
	"sync"
	"time"

	"github.com/grillinr/nq/db"
	lru "github.com/hashicorp/golang-lru/v2"
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
	Repo               db.Repository
	searchStatuses     *lru.Cache[string, searchStatus]
	searchMu           sync.RWMutex
	stopCleanup        chan struct{}
	cleanupDone        chan struct{}
	mediaCreationMu    sync.Mutex
	mediaCreationLocks map[string]*sync.Mutex // key: "type:title:year"
}

// NewResolver creates a new resolver with database repository
func NewResolver(repo db.Repository) *Resolver {
	cache, err := lru.New[string, searchStatus](maxSearchStatuses)
	if err != nil {
		panic(err) // Should never happen with valid size
	}

	resolver := &Resolver{
		Repo:               repo,
		searchStatuses:     cache,
		stopCleanup:        make(chan struct{}),
		cleanupDone:        make(chan struct{}),
		mediaCreationLocks: make(map[string]*sync.Mutex),
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

// acquireMediaLock gets or creates a lock for a specific media item to prevent concurrent creation.
// Returns the lock and a release function that should be called with defer.
func (r *Resolver) acquireMediaLock(mediaType, title string, year int) (*sync.Mutex, func()) {
	key := fmt.Sprintf("%s:%s:%d", mediaType, title, year)

	r.mediaCreationMu.Lock()
	lock, exists := r.mediaCreationLocks[key]
	if !exists {
		lock = &sync.Mutex{}
		r.mediaCreationLocks[key] = lock
	}
	r.mediaCreationMu.Unlock()

	lock.Lock()

	// Return a cleanup function that releases the lock and potentially removes it from the map
	return lock, func() {
		lock.Unlock()

		// Clean up the lock from the map if it's no longer needed
		// This prevents the map from growing indefinitely
		r.mediaCreationMu.Lock()
		defer r.mediaCreationMu.Unlock()

		// Only remove if no other goroutine is waiting on this lock
		// We do a tryLock to check - if we can acquire it, no one else is waiting
		if lock.TryLock() {
			delete(r.mediaCreationLocks, key)
			lock.Unlock()
		}
	}
}
