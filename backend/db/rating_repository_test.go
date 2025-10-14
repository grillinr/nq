package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

func TestRatingRepository_CreateRating(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and media first
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	// Test data
	score := 8.5

	// Test creating a rating
	rating, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
	if err != nil {
		t.Fatalf("Failed to create rating: %v", err)
	}

	// Verify rating properties
	if rating.User == nil || rating.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, rating.User.ID)
	}
	if rating.Media == nil || rating.Media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, rating.Media.GetID())
	}
	if rating.Score != score {
		t.Errorf("Expected score %f, got %f", score, rating.Score)
	}
	if rating.RatedAt == "" {
		t.Error("Expected non-empty rating timestamp")
	}
}

func TestRatingRepository_GetRating(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)
	score := 7.0

	createdRating, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
	if err != nil {
		t.Fatalf("Failed to create test rating: %v", err)
	}

	// Test getting rating
	rating, err := repo.GetRating(ctx, user.ID, movie.ID)
	if err != nil {
		t.Fatalf("Failed to get rating: %v", err)
	}

	// Verify rating properties
	if rating.User == nil || rating.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, rating.User.ID)
	}
	if rating.Media == nil || rating.Media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, rating.Media.GetID())
	}
	if rating.Score != createdRating.Score {
		t.Errorf("Expected score %f, got %f", createdRating.Score, rating.Score)
	}

	// Test getting non-existent rating
	anotherUser := createTestUser(t, repo, ctx)
	_, err = repo.GetRating(ctx, anotherUser.ID, movie.ID)
	if err == nil {
		t.Error("Expected error when getting non-existent rating")
	}
}

func TestRatingRepository_GetUserRatings(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and multiple ratings
	user := createTestUser(t, repo, ctx)

	ratings := make([]*model.Rating, 3)
	for i := 0; i < 3; i++ {
		movie := createTestMovie(t, repo, ctx)
		score := float64(5 + i) // 5.0, 6.0, 7.0

		rating, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
		if err != nil {
			t.Fatalf("Failed to create test rating %d: %v", i, err)
		}
		ratings[i] = rating
	}

	// Test getting user ratings
	userRatings, err := repo.GetUserRatings(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user ratings: %v", err)
	}

	// Should have at least our test ratings
	if len(userRatings) < 3 {
		t.Errorf("Expected at least 3 ratings, got %d", len(userRatings))
	}

	// Verify our test ratings are in the results by checking user and media IDs
	foundRatings := 0
	for _, createdRating := range ratings {
		for _, fetchedRating := range userRatings {
			if fetchedRating.User != nil && fetchedRating.User.ID == createdRating.User.ID &&
				fetchedRating.Media != nil && fetchedRating.Media.GetID() == createdRating.Media.GetID() {
				foundRatings++
				break
			}
		}
	}

	if foundRatings != 3 {
		t.Errorf("Expected to find all 3 created ratings, found %d", foundRatings)
	}
}

func TestRatingRepository_GetMediaRatings(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test movie and multiple ratings for it
	movie := createTestMovie(t, repo, ctx)

	ratings := make([]*model.Rating, 3)
	for i := 0; i < 3; i++ {
		user := createTestUser(t, repo, ctx)
		score := float64(6 + i) // 6.0, 7.0, 8.0

		rating, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
		if err != nil {
			t.Fatalf("Failed to create test rating %d: %v", i, err)
		}
		ratings[i] = rating
	}

	// Test getting media ratings
	mediaRatings, err := repo.GetMediaRatings(ctx, movie.ID)
	if err != nil {
		t.Fatalf("Failed to get media ratings: %v", err)
	}

	// Should have at least our test ratings
	if len(mediaRatings) < 3 {
		t.Errorf("Expected at least 3 ratings, got %d", len(mediaRatings))
	}

	// Verify our test ratings are in the results by checking user and media IDs
	foundRatings := 0
	for _, createdRating := range ratings {
		for _, fetchedRating := range mediaRatings {
			if fetchedRating.User != nil && fetchedRating.User.ID == createdRating.User.ID &&
				fetchedRating.Media != nil && fetchedRating.Media.GetID() == createdRating.Media.GetID() {
				foundRatings++
				break
			}
		}
	}

	if foundRatings != 3 {
		t.Errorf("Expected to find all 3 created ratings, found %d", foundRatings)
	}
}

func TestRatingRepository_UpdateRating(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)
	originalScore := 6.0

	_, err := repo.CreateRating(ctx, user.ID, movie.ID, originalScore)
	if err != nil {
		t.Fatalf("Failed to create test rating: %v", err)
	}

	// Test updating rating
	newScore := 9.0
	updatedRating, err := repo.UpdateRating(ctx, user.ID, movie.ID, newScore)
	if err != nil {
		t.Fatalf("Failed to update rating: %v", err)
	}

	// Verify updated properties
	if updatedRating.Score != newScore {
		t.Errorf("Expected score %f, got %f", newScore, updatedRating.Score)
	}
	if updatedRating.User == nil || updatedRating.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, updatedRating.User.ID)
	}
	if updatedRating.Media == nil || updatedRating.Media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, updatedRating.Media.GetID())
	}

	// Test updating non-existent rating
	anotherUser := createTestUser(t, repo, ctx)
	_, err = repo.UpdateRating(ctx, anotherUser.ID, movie.ID, newScore)
	if err == nil {
		t.Error("Expected error when updating non-existent rating")
	}
}

func TestRatingRepository_DeleteRating(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)
	score := 7.5

	_, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
	if err != nil {
		t.Fatalf("Failed to create test rating: %v", err)
	}

	// Test deleting rating
	err = repo.DeleteRating(ctx, user.ID, movie.ID)
	if err != nil {
		t.Fatalf("Failed to delete rating: %v", err)
	}

	// Verify rating is deleted
	_, err = repo.GetRating(ctx, user.ID, movie.ID)
	if err == nil {
		t.Error("Expected error when getting deleted rating")
	}

	// Test deleting non-existent rating
	err = repo.DeleteRating(ctx, user.ID, movie.ID)
	if err == nil {
		t.Error("Expected error when deleting non-existent rating")
	}
}

func TestRatingRepository_GetAverageRating(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test movie and multiple ratings for it
	movie := createTestMovie(t, repo, ctx)

	scores := []float64{6.0, 8.0, 10.0} // Average should be 8.0
	for i, score := range scores {
		user := createTestUser(t, repo, ctx)
		_, err := repo.CreateRating(ctx, user.ID, movie.ID, score)
		if err != nil {
			t.Fatalf("Failed to create test rating %d: %v", i, err)
		}
	}

	// Test getting average rating
	avgRating, err := repo.GetAverageRating(ctx, movie.ID)
	if err != nil {
		t.Fatalf("Failed to get average rating: %v", err)
	}

	if avgRating == nil {
		t.Fatal("Expected non-nil average rating")
	}

	expectedAvg := 8.0
	if *avgRating != expectedAvg {
		t.Errorf("Expected average rating %f, got %f", expectedAvg, *avgRating)
	}

	// Test getting average rating for media with no ratings
	anotherMovie := createTestMovie(t, repo, ctx)
	avgRating, err = repo.GetAverageRating(ctx, anotherMovie.ID)
	if err != nil {
		t.Fatalf("Failed to get average rating for unrated media: %v", err)
	}

	if avgRating != nil {
		t.Errorf("Expected nil average rating for unrated media, got %v", avgRating)
	}

	// Test getting average rating for non-existent media
	nonExistentID := uuid.New()
	avgRating, err = repo.GetAverageRating(ctx, nonExistentID)
	if err != nil {
		t.Fatalf("Failed to get average rating for non-existent media: %v", err)
	}

	if avgRating != nil {
		t.Errorf("Expected nil average rating for non-existent media, got %v", avgRating)
	}
}
