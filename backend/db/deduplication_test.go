package db

import (
	"context"
	"testing"

	"github.com/grillinr/nq/graph/model"
)

func TestDeduplication_PersonAndProductionCompany(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Use variants of the same names
	cast1 := []string{"José González"}
	pcs1 := []string{"Acme, Inc."}

	input1 := model.CreateMovieInput{
		Title:               GenerateTestTitle("movie"),
		Cast:                cast1,
		ProductionCompanies: pcs1,
	}

	m1, err := repo.CreateMovie(ctx, input1)
	if err != nil {
		t.Fatalf("failed to create first movie: %v", err)
	}

	// Second movie with slightly different variants
	cast2 := []string{"Jose Gonzalez"}
	pcs2 := []string{"Acme Inc"}

	input2 := model.CreateMovieInput{
		Title:               GenerateTestTitle("movie"),
		Cast:                cast2,
		ProductionCompanies: pcs2,
	}

	m2, err := repo.CreateMovie(ctx, input2)
	if err != nil {
		t.Fatalf("failed to create second movie: %v", err)
	}

	// Fetch movies
	f1, err := repo.GetMovieByID(ctx, m1.ID)
	if err != nil {
		t.Fatalf("failed to fetch movie1: %v", err)
	}
	f2, err := repo.GetMovieByID(ctx, m2.ID)
	if err != nil {
		t.Fatalf("failed to fetch movie2: %v", err)
	}

	if len(f1.Cast) == 0 || len(f2.Cast) == 0 {
		t.Fatalf("expected cast to be present")
	}
	// If deduplication worked, the normalized names should be equal and refer to same person name
	if f1.Cast[0].Name != f2.Cast[0].Name {
		t.Errorf("expected normalized cast names to match: %q vs %q", f1.Cast[0].Name, f2.Cast[0].Name)
	}

	// Production company names
	if len(f1.ProductionCompanies) == 0 || len(f2.ProductionCompanies) == 0 {
		t.Fatalf("expected production companies to be present")
	}
	if f1.ProductionCompanies[0].Name != f2.ProductionCompanies[0].Name {
		t.Errorf("expected normalized production company names to match: %q vs %q", f1.ProductionCompanies[0].Name, f2.ProductionCompanies[0].Name)
	}
}
