package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

func TestUserRepository_CreateUser(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("google"),
	}

	// Test creating a user
	user, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user properties
	if user.Name != input.Name {
		t.Errorf("Expected name %s, got %s", input.Name, user.Name)
	}
	if user.Email != input.Email {
		t.Errorf("Expected email %s, got %s", input.Email, user.Email)
	}
	if user.AuthProvider == nil || *user.AuthProvider != *input.AuthProvider {
		t.Errorf("Expected auth provider %v, got %v", input.AuthProvider, user.AuthProvider)
	}
	if user.ID == uuid.Nil {
		t.Error("Expected non-nil user ID")
	}

	// Test duplicate email should fail
	_, err = repo.CreateUser(ctx, input)
	if err == nil {
		t.Error("Expected error when creating user with duplicate email")
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test user first
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("github"),
	}

	createdUser, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test getting user by ID
	user, err := repo.GetUserByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	// Verify user properties
	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %v, got %v", createdUser.ID, user.ID)
	}
	if user.Name != createdUser.Name {
		t.Errorf("Expected name %s, got %s", createdUser.Name, user.Name)
	}
	if user.Email != createdUser.Email {
		t.Errorf("Expected email %s, got %s", createdUser.Email, user.Email)
	}

	// Test non-existent user
	nonExistentID := uuid.New()
	_, err = repo.GetUserByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test user first
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("apple"),
	}

	createdUser, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test getting user by email
	user, err := repo.GetUserByEmail(ctx, createdUser.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	// Verify user properties
	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %v, got %v", createdUser.ID, user.ID)
	}
	if user.Email != createdUser.Email {
		t.Errorf("Expected email %s, got %s", createdUser.Email, user.Email)
	}

	// Test non-existent email
	_, err = repo.GetUserByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("Expected error when getting user with non-existent email")
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create multiple test users
	users := make([]*model.User, 3)
	for i := 0; i < 3; i++ {
		input := model.CreateUserInput{
			Name:         GenerateTestName(),
			Email:        GenerateTestEmail(),
			AuthProvider: stringPointer("test"),
		}

		user, err := repo.CreateUser(ctx, input)
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", i, err)
		}
		users[i] = user
	}

	// Test getting all users
	allUsers, err := repo.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	// Should have at least our test users
	if len(allUsers) < 3 {
		t.Errorf("Expected at least 3 users, got %d", len(allUsers))
	}

	// Verify our test users are in the results
	foundUsers := 0
	for _, createdUser := range users {
		for _, fetchedUser := range allUsers {
			if fetchedUser.ID == createdUser.ID {
				foundUsers++
				break
			}
		}
	}

	if foundUsers != 3 {
		t.Errorf("Expected to find all 3 created users, found %d", foundUsers)
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test user first
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("original"),
	}

	createdUser, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test updating user
	newName := GenerateTestName()
	newEmail := GenerateTestEmail()
	updateInput := model.UpdateUserInput{
		Name:  stringPointer(newName),
		Email: stringPointer(newEmail),
	}

	updatedUser, err := repo.UpdateUser(ctx, createdUser.ID, updateInput)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify updated properties
	if updatedUser.Name != *updateInput.Name {
		t.Errorf("Expected name %s, got %s", *updateInput.Name, updatedUser.Name)
	}
	if updatedUser.Email != *updateInput.Email {
		t.Errorf("Expected email %s, got %s", *updateInput.Email, updatedUser.Email)
	}
	// AuthProvider should remain unchanged
	if updatedUser.AuthProvider == nil || *updatedUser.AuthProvider != *createdUser.AuthProvider {
		t.Errorf("Expected auth provider to remain %v, got %v", createdUser.AuthProvider, updatedUser.AuthProvider)
	}

	// Test updating non-existent user
	nonExistentID := uuid.New()
	_, err = repo.UpdateUser(ctx, nonExistentID, updateInput)
	if err == nil {
		t.Error("Expected error when updating non-existent user")
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test user first
	input := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("test"),
	}

	createdUser, err := repo.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test deleting user
	err = repo.DeleteUser(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user is deleted
	_, err = repo.GetUserByID(ctx, createdUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}

	// Test deleting non-existent user
	nonExistentID := uuid.New()
	err = repo.DeleteUser(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}
}

func TestUserRepository_AddToFavorites(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test user
	userInput := model.CreateUserInput{
		Name:         GenerateTestName(),
		Email:        GenerateTestEmail(),
		AuthProvider: stringPointer("test"),
	}

	user, err := repo.CreateUser(ctx, userInput)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a test movie
	movieInput := model.CreateMovieInput{
		Title:       GenerateTestTitle("movie"),
		Description: stringPointer("Test movie description"),
	}

	movie, err := repo.CreateMovie(ctx, movieInput)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	// Test adding to favorites
	err = repo.AddToFavorites(ctx, user.ID, movie.ID)
	if err != nil {
		t.Fatalf("Failed to add to favorites: %v", err)
	}

	// Test adding duplicate favorite (should handle gracefully)
	err = repo.AddToFavorites(ctx, user.ID, movie.ID)
	if err != nil {
		t.Logf("Adding duplicate favorite returned error: %v", err)
	}

	// Test adding favorite with non-existent user
	nonExistentUserID := uuid.New()
	err = repo.AddToFavorites(ctx, nonExistentUserID, movie.ID)
	if err == nil {
		t.Error("Expected error when adding favorite with non-existent user")
	}

	// Test adding favorite with non-existent media
	nonExistentMediaID := uuid.New()
	err = repo.AddToFavorites(ctx, user.ID, nonExistentMediaID)
	if err == nil {
		t.Error("Expected error when adding favorite with non-existent media")
	}
}

// Helper function for string pointers
func stringPointer(s string) *string {
	return &s
}
