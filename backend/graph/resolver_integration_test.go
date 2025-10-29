package graph

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/grillinr/nq/db"
	"github.com/grillinr/nq/graph/model"
)

func loadEnvUpwards(filename string, maxDepth int) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	p := cwd
	for i := 0; i <= maxDepth; i++ {
		candidate := filepath.Join(p, filename)
		if _, statErr := os.Stat(candidate); statErr == nil {
			// found
			_ = godotenv.Load(candidate)
			return nil
		}
		p = filepath.Dir(p)
	}
	return fmt.Errorf("%s not found in cwd or %d parent directories", filename, maxDepth)
}

func TestResolver_CreateMovieMutation_E2E(t *testing.T) {
	// Connect to real database (skip if unavailable)
	dbConn, err := db.NewDatabase()
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	repo := db.NewNeo4jRepository(dbConn)
	resolver := NewResolver(repo)

	ctx := context.Background()

	castNames := []string{"integration-cast-1", "integration-cast-2"}
	crewNames := []string{"integration-crew-1", "integration-crew-2"}

	input := model.CreateMovieInput{
		Title: "integration-test-movie",
		Cast:  castNames,
		Crew:  crewNames,
	}

	created, err := resolver.Mutation().CreateMovie(ctx, input)
	if err != nil {
		t.Fatalf("Resolver CreateMovie failed: %v", err)
	}
	if created == nil {
		t.Fatalf("Resolver returned nil movie")
	}
	if created.ID == uuid.Nil {
		t.Fatalf("Expected non-nil movie ID from resolver")
	}

	// Verify persisted
	fetched, err := repo.GetMovieByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to fetch movie by ID: %v", err)
	}
	if len(fetched.Cast) != len(castNames) {
		t.Errorf("Expected %d cast members, got %d", len(castNames), len(fetched.Cast))
	}
	if len(fetched.Crew) != len(crewNames) {
		t.Errorf("Expected %d crew members, got %d", len(crewNames), len(fetched.Crew))
	}

	// Cleanup - delete the created movie node and related persons
	_, _ = dbConn.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) { return nil, nil })
	// Fallback cleanup by ID
	_, _ = dbConn.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (m:Movie {id: $id}) DETACH DELETE m"
		_, err := tx.Run(ctx, query, map[string]any{"id": created.ID.String()})
		return nil, err
	})
}
