package graph

import (
	"time"
)

func (r *Resolver) setSearchStatus(id string, state searchState, completedAt *time.Time) {
	r.searchMu.Lock()
	defer r.searchMu.Unlock()

	r.searchStatuses[id] = searchStatus{
		state:       state,
		completedAt: completedAt,
	}
}

func (r *Resolver) getSearchStatus(id string) searchStatus {
	r.searchMu.RLock()
	defer r.searchMu.RUnlock()

	if status, ok := r.searchStatuses[id]; ok {
		return status
	}
	return searchStatus{state: searchStateIdle}
}
