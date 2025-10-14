package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

func TestRecommendationRepository_CreateRecommendation(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and media first
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	// Test data
	source := "manual"
	score := 8.5

	// Test creating a recommendation
	recommendation, err := repo.CreateRecommendation(ctx, user.ID, movie.ID, nil, &source, &score)
	if err != nil {
		t.Fatalf("Failed to create recommendation: %v", err)
	}

	// Verify recommendation properties
	if recommendation.User == nil || recommendation.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, recommendation.User.ID)
	}
	if recommendation.Media == nil || recommendation.Media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, recommendation.Media.GetID())
	}
	if recommendation.Source == nil || *recommendation.Source != source {
		t.Errorf("Expected source %s, got %v", source, recommendation.Source)
	}
	if recommendation.Score == nil || *recommendation.Score != score {
		t.Errorf("Expected score %f, got %v", score, recommendation.Score)
	}
	if recommendation.ID == uuid.Nil {
		t.Error("Expected non-nil recommendation ID")
	}
}

func TestRecommendationRepository_GetRecommendationByID(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)
	source := "test"
	score := 7.5

	createdRecommendation, err := repo.CreateRecommendation(ctx, user.ID, movie.ID, nil, &source, &score)
	if err != nil {
		t.Fatalf("Failed to create test recommendation: %v", err)
	}

	// Test getting recommendation
	recommendation, err := repo.GetRecommendationByID(ctx, createdRecommendation.ID)
	if err != nil {
		t.Fatalf("Failed to get recommendation: %v", err)
	}

	// Verify recommendation properties
	if recommendation.ID != createdRecommendation.ID {
		t.Errorf("Expected ID %v, got %v", createdRecommendation.ID, recommendation.ID)
	}
	if recommendation.Source == nil || *recommendation.Source != source {
		t.Errorf("Expected source %s, got %v", source, recommendation.Source)
	}
	if recommendation.User == nil || recommendation.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, recommendation.User.ID)
	}
	if recommendation.Media == nil || recommendation.Media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, recommendation.Media.GetID())
	}

	// Test getting non-existent recommendation
	nonExistentID := uuid.New()
	_, err = repo.GetRecommendationByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent recommendation")
	}
}

func TestRecommendationRepository_GetRecommendations(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and multiple recommendations
	user := createTestUser(t, repo, ctx)

	recommendations := make([]*model.Recommendation, 3)
	for i := 0; i < 3; i++ {
		movie := createTestMovie(t, repo, ctx)
		source := fmt.Sprintf("source_%d", i+1)
		score := float64(5 + i) // 5.0, 6.0, 7.0

		recommendation, err := repo.CreateRecommendation(ctx, user.ID, movie.ID, nil, &source, &score)
		if err != nil {
			t.Fatalf("Failed to create test recommendation %d: %v", i, err)
		}
		recommendations[i] = recommendation
	}

	// Test getting user recommendations
	userRecommendations, err := repo.GetRecommendations(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user recommendations: %v", err)
	}

	// Should have at least our test recommendations
	if len(userRecommendations) < 3 {
		t.Errorf("Expected at least 3 recommendations, got %d", len(userRecommendations))
	}

	// Verify our test recommendations are in the results
	foundRecommendations := 0
	for _, createdRecommendation := range recommendations {
		for _, fetchedRecommendation := range userRecommendations {
			if fetchedRecommendation.ID == createdRecommendation.ID {
				foundRecommendations++
				break
			}
		}
	}

	if foundRecommendations != 3 {
		t.Errorf("Expected to find all 3 created recommendations, found %d", foundRecommendations)
	}
}

func TestRecommendationRepository_DeleteRecommendation(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)
	source := "test"
	score := 8.0

	createdRecommendation, err := repo.CreateRecommendation(ctx, user.ID, movie.ID, nil, &source, &score)
	if err != nil {
		t.Fatalf("Failed to create test recommendation: %v", err)
	}

	// Test deleting recommendation
	err = repo.DeleteRecommendation(ctx, createdRecommendation.ID)
	if err != nil {
		t.Fatalf("Failed to delete recommendation: %v", err)
	}

	// Verify recommendation is deleted
	_, err = repo.GetRecommendationByID(ctx, createdRecommendation.ID)
	if err == nil {
		t.Error("Expected error when getting deleted recommendation")
	}

	// Test deleting non-existent recommendation
	err = repo.DeleteRecommendation(ctx, createdRecommendation.ID)
	if err == nil {
		t.Error("Expected error when deleting non-existent recommendation")
	}
}
