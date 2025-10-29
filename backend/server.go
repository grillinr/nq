package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grillinr/nq/db"
	"github.com/grillinr/nq/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
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

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return recoverMiddleware(srv)
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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", NewGraphQLHandler(repo))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
