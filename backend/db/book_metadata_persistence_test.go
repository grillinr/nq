package db

import (
	"context"
	"testing"

	"github.com/grillinr/nq/graph/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestBookMetadataPersistence(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	input := model.CreateBookInput{
		Title:      GenerateTestTitle("book"),
		Authors:    []string{"test-author-1", "test-author-2"},
		Publishers: []string{"publisher-a", "publisher-b"},
		Subjects:   []string{"Subject X", "Subject Y"},
	}

	created, err := repo.CreateBook(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create book with metadata: %v", err)
	}

	fetched, err := repo.GetBookByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to fetch book by ID: %v", err)
	}

	// Verify publishers
	if len(fetched.Publishers) != len(input.Publishers) {
		t.Fatalf("Expected %d publishers, got %d", len(input.Publishers), len(fetched.Publishers))
	}

	// Verify publisher nodes exist
	cntRes, err := testDB.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (p:Publisher) WHERE p.name IN $names RETURN count(p) as cnt"
		params := map[string]any{"names": input.Publishers}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			rec := result.Record()
			val, _ := rec.Get("cnt")
			return val, nil
		}
		return int64(0), nil
	})
	if err != nil {
		t.Fatalf("Failed to query publisher nodes: %v", err)
	}
	cnt, ok := cntRes.(int64)
	if !ok {
		t.Fatalf("Unexpected count type: %T", cntRes)
	}
	if cnt < int64(len(input.Publishers)) {
		t.Fatalf("Expected at least %d publisher nodes, found %d", len(input.Publishers), cnt)
	}

	// Verify PUBLISHED relationships
	relRes, err := testDB.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (p:Publisher)-[r:PUBLISHED]->(b:Book {id: $id}) RETURN count(r) as cnt"
		params := map[string]any{"id": created.ID.String()}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			rec := result.Record()
			val, _ := rec.Get("cnt")
			return val, nil
		}
		return int64(0), nil
	})
	if err != nil {
		t.Fatalf("Failed to query PUBLISHED relationships: %v", err)
	}
	relCnt, ok := relRes.(int64)
	if !ok {
		t.Fatalf("Unexpected rel count type: %T", relRes)
	}
	if relCnt < int64(len(input.Publishers)) {
		t.Fatalf("Expected at least %d PUBLISHED relationships, found %d", len(input.Publishers), relCnt)
	}

	// Verify authors
	if len(fetched.Authors) != len(input.Authors) {
		t.Fatalf("Expected %d authors, got %d", len(input.Authors), len(fetched.Authors))
	}

	// Verify subjects
	if len(fetched.Subjects) != len(input.Subjects) {
		t.Fatalf("Expected %d subjects, got %d", len(input.Subjects), len(fetched.Subjects))
	}
}
