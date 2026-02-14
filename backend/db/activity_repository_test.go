package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

func TestActivityRepository_CreateActivity(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and media first
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	// Test data
	input := model.CreateActivityInput{
		MediaID:    movie.ID,
		StatusID:   1, // Assuming 1 is a valid status
		Rating:     float64Pointer(8.5),
		Review:     stringPointer("Great movie!"),
		StartedAt:  stringPointer("2023-01-01T00:00:00Z"),
		FinishedAt: stringPointer("2023-01-02T00:00:00Z"),
	}

	// Test creating an activity
	activity, err := repo.CreateActivity(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	// Verify activity properties
	if activity.User == nil || activity.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, activity.User.ID)
	}
	if activity.Media == nil || activity.Media.GetID() != input.MediaID {
		t.Errorf("Expected media ID %v, got %v", input.MediaID, activity.Media.GetID())
	}
	if activity.Status == nil || activity.Status.ID != input.StatusID {
		t.Errorf("Expected status ID %d, got %d", input.StatusID, activity.Status.ID)
	}
	if activity.Rating == nil || *activity.Rating != *input.Rating {
		t.Errorf("Expected rating %v, got %v", input.Rating, activity.Rating)
	}
	if activity.ID == uuid.Nil {
		t.Error("Expected non-nil activity ID")
	}
}

func TestActivityRepository_GetActivityByID(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	input := model.CreateActivityInput{
		MediaID:  movie.ID,
		StatusID: 2,
		Rating:   float64Pointer(7.0),
	}

	createdActivity, err := repo.CreateActivity(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}

	// Test getting activity by ID
	activity, err := repo.GetActivityByID(ctx, createdActivity.ID)
	if err != nil {
		t.Fatalf("Failed to get activity by ID: %v", err)
	}

	// Verify activity properties
	if activity.ID != createdActivity.ID {
		t.Errorf("Expected ID %v, got %v", createdActivity.ID, activity.ID)
	}
	if activity.User == nil || createdActivity.User == nil || activity.User.ID != createdActivity.User.ID {
		t.Errorf("Expected user ID %v, got %v", createdActivity.User.ID, activity.User.ID)
	}
	if activity.Media == nil || createdActivity.Media == nil || activity.Media.GetID() != createdActivity.Media.GetID() {
		t.Errorf("Expected media ID %v, got %v", createdActivity.Media.GetID(), activity.Media.GetID())
	}

	// Test non-existent activity
	nonExistentID := uuid.New()
	_, err = repo.GetActivityByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent activity")
	}
}

func TestActivityRepository_GetUserActivities(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test user and multiple activities
	user := createTestUser(t, repo, ctx)

	activities := make([]*model.UserActivity, 3)
	for i := 0; i < 3; i++ {
		movie := createTestMovie(t, repo, ctx)
		input := model.CreateActivityInput{
			MediaID:  movie.ID,
			StatusID: int32(i + 1),
		}

		activity, err := repo.CreateActivity(ctx, user.ID, input)
		if err != nil {
			t.Fatalf("Failed to create test activity %d: %v", i, err)
		}
		activities[i] = activity
	}

	// Test getting user activities
	userActivities, err := repo.GetUserActivities(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user activities: %v", err)
	}

	// Should have at least our test activities
	if len(userActivities) < 3 {
		t.Errorf("Expected at least 3 activities, got %d", len(userActivities))
	}

	// Verify our test activities are in the results
	foundActivities := 0
	for _, createdActivity := range activities {
		for _, fetchedActivity := range userActivities {
			if fetchedActivity.ID == createdActivity.ID {
				foundActivities++
				break
			}
		}
	}

	if foundActivities != 3 {
		t.Errorf("Expected to find all 3 created activities, found %d", foundActivities)
	}

	// Test getting activities for non-existent user
	nonExistentUserID := uuid.New()
	emptyActivities, err := repo.GetUserActivities(ctx, nonExistentUserID)
	if err != nil {
		t.Fatalf("Failed to get activities for non-existent user: %v", err)
	}
	if len(emptyActivities) != 0 {
		t.Errorf("Expected 0 activities for non-existent user, got %d", len(emptyActivities))
	}
}

func TestActivityRepository_GetMediaActivities(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test movie and multiple user activities for it
	movie := createTestMovie(t, repo, ctx)

	activities := make([]*model.UserActivity, 3)
	for i := 0; i < 3; i++ {
		user := createTestUser(t, repo, ctx)
		input := model.CreateActivityInput{
			MediaID:  movie.ID,
			StatusID: int32(i + 1),
		}

		activity, err := repo.CreateActivity(ctx, user.ID, input)
		if err != nil {
			t.Fatalf("Failed to create test activity %d: %v", i, err)
		}
		activities[i] = activity
	}

	// Test getting media activities
	mediaActivities, err := repo.GetMediaActivities(ctx, movie.ID)
	if err != nil {
		t.Fatalf("Failed to get media activities: %v", err)
	}

	// Should have at least our test activities
	if len(mediaActivities) < 3 {
		t.Errorf("Expected at least 3 activities, got %d", len(mediaActivities))
	}

	// Verify our test activities are in the results
	foundActivities := 0
	for _, createdActivity := range activities {
		for _, fetchedActivity := range mediaActivities {
			if fetchedActivity.ID == createdActivity.ID {
				foundActivities++
				break
			}
		}
	}

	if foundActivities != 3 {
		t.Errorf("Expected to find all 3 created activities, found %d", foundActivities)
	}
}

func TestActivityRepository_UpdateActivity(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	input := model.CreateActivityInput{
		MediaID:  movie.ID,
		StatusID: 1,
		Rating:   float64Pointer(6.0),
		Review:   stringPointer("Original review"),
	}

	createdActivity, err := repo.CreateActivity(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}

	// Test updating activity
	newStatusID := int32(2)
	newRating := 8.5
	newReview := "Updated review"
	newFinishedAt := "2023-12-31T23:59:59Z"

	updateInput := model.UpdateActivityInput{
		StatusID:   &newStatusID,
		Rating:     &newRating,
		Review:     &newReview,
		FinishedAt: &newFinishedAt,
	}

	updatedActivity, err := repo.UpdateActivity(ctx, user.ID, createdActivity.ID, updateInput)
	if err != nil {
		t.Fatalf("Failed to update activity: %v", err)
	}

	// Verify updated properties
	if updatedActivity.Status == nil || updatedActivity.Status.ID != newStatusID {
		t.Errorf("Expected status ID %d, got %d", newStatusID, updatedActivity.Status.ID)
	}
	if updatedActivity.Rating == nil || *updatedActivity.Rating != newRating {
		t.Errorf("Expected rating %f, got %v", newRating, updatedActivity.Rating)
	}
	if updatedActivity.Review == nil || *updatedActivity.Review != newReview {
		t.Errorf("Expected review %s, got %v", newReview, updatedActivity.Review)
	}

	// Test updating non-existent activity
	nonExistentID := uuid.New()
	nonExistentInput := model.UpdateActivityInput{
		StatusID: &newStatusID,
	}
	_, err = repo.UpdateActivity(ctx, user.ID, nonExistentID, nonExistentInput)
	if err == nil {
		t.Error("Expected error when updating non-existent activity")
	}
}

func TestActivityRepository_DeleteActivity(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create test data
	user := createTestUser(t, repo, ctx)
	movie := createTestMovie(t, repo, ctx)

	input := model.CreateActivityInput{
		MediaID:  movie.ID,
		StatusID: 1,
	}

	createdActivity, err := repo.CreateActivity(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}

	// Test deleting activity
	err = repo.DeleteActivity(ctx, createdActivity.ID)
	if err != nil {
		t.Fatalf("Failed to delete activity: %v", err)
	}

	// Verify activity is deleted
	_, err = repo.GetActivityByID(ctx, createdActivity.ID)
	if err == nil {
		t.Error("Expected error when getting deleted activity")
	}

	// Test deleting non-existent activity
	nonExistentID := uuid.New()
	err = repo.DeleteActivity(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when deleting non-existent activity")
	}
}

// Helper functions to create test data
func createTestUser(t *testing.T, repo *Neo4jRepository, ctx context.Context) *model.User {
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("test"),
	}

	user, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestMovie(t *testing.T, repo *Neo4jRepository, ctx context.Context) *model.Movie {
	input := model.CreateMovieInput{
		Title:       GenerateTestTitle("movie"),
		Description: stringPointer("Test movie"),
	}

	movie, err := repo.CreateMovie(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}
	return movie
}
