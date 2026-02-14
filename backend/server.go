package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grillinr/nq/auth"
	"github.com/grillinr/nq/db"
	"github.com/grillinr/nq/graph"
	"github.com/grillinr/nq/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/time/rate"
)

const defaultPort = "8080"

// recoverMiddleware returns an http.Handler that recovers from panics and returns HTTP 500
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("recovered from panic: %v", rec)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NewGraphQLHandler creates and returns an http.Handler for the GraphQL endpoint.
// If repo is nil, a nil repository is passed to the resolver — resolvers must handle a nil repo if used.
func NewGraphQLHandler(repo db.Repository) http.Handler {
	// Create resolver with repository (repo may be nil in tests)
	resolver := graph.NewResolver(repo)

	// Create GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Only enable introspection in development
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		srv.Use(extension.Introspection{})
	}

	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Add query complexity limit to prevent abuse
	srv.Use(extension.FixedComplexityLimit(1000))

	return http.TimeoutHandler(recoverMiddleware(srv), 5*time.Minute, "Request timeout")
}

// corsMiddleware sets CORS headers based on allowed origins from environment.
// It handles preflight OPTIONS requests and forwards other requests to next.
// Deprecated: Use middleware.CORS instead
func corsMiddleware(next http.Handler) http.Handler {
	return middleware.CORS(next)
}

func GraphQL() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize database
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize database constraints and indexes
	ctx := context.Background()
	if err := database.InitializeDatabase(ctx); err != nil {
		log.Printf("Warning: Failed to initialize database constraints: %v", err)
	}

	// Create repository
	repo := db.NewNeo4jRepository(database)

	validator, err := auth.NewValidatorFromEnv()
	if err != nil {
		log.Printf("Warning: auth disabled: %v", err)
	}

	// Setup middleware chain with security features
	rateLimiter := middleware.NewRateLimiter(rate.Limit(100), 200) // 100 req/sec, burst 200

	// GraphQL playground - only in development
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	}

	// Build middleware chain: security headers -> rate limit -> CORS -> auth -> handler
	graphqlHandler := NewGraphQLHandler(repo)
	graphqlHandler = auth.AuthMiddleware(validator, repo)(graphqlHandler)
	graphqlHandler = corsMiddleware(graphqlHandler)
	graphqlHandler = rateLimiter.Limit(graphqlHandler)
	graphqlHandler = middleware.SecurityHeaders(graphqlHandler)

	http.Handle("/graphql", graphqlHandler)

	// Start server with optional TLS
	enableTLS := os.Getenv("ENABLE_TLS") == "true"
	if enableTLS {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")

		if certFile == "" || keyFile == "" {
			log.Fatal("TLS enabled but TLS_CERT_FILE or TLS_KEY_FILE not set")
		}

		log.Printf("connect to https://localhost:%s/ for GraphQL playground", port)
		log.Fatal(http.ListenAndServeTLS(":"+port, certFile, keyFile, nil))
	} else {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Println("WARNING: TLS is disabled. Set ENABLE_TLS=true for production")
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
