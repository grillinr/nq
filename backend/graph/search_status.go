package graph

import "time"

type searchStatus struct {
	state       searchState
	completedAt *time.Time
}

type searchState string

const (
	searchStateIdle      searchState = "IDLE"
	searchStateRunning   searchState = "RUNNING"
	searchStateCompleted searchState = "COMPLETED"
)
