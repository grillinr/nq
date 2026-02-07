package config

import "time"

// HTTP client timeout durations
const (
	// AuthHTTPTimeout is the timeout for auth-related HTTP requests (userinfo endpoint)
	AuthHTTPTimeout = 5 * time.Second

	// IntegrationHTTPTimeout is the timeout for external integration API calls
	IntegrationHTTPTimeout = 30 * time.Second

	// MetadataHTTPTimeout is the timeout for metadata fetching API calls
	MetadataHTTPTimeout = 10 * time.Second

	// ServerRequestTimeout is the maximum time allowed for a single request
	ServerRequestTimeout = 5 * time.Minute
)

// Cache and refresh intervals
const (
	// JWKSCacheDuration is how long to cache JWKS keys before refreshing
	JWKSCacheDuration = 10 * time.Minute
)

// Search and connection limits
const (
	// DefaultMaxConnections is the default maximum number of connections to search
	DefaultMaxConnections = 25
)
