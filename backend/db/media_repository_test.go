package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

func TestMediaRepository_CreateMovie(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateMovieInput{
		Title:       GenerateTestTitle("movie"),
		ReleaseDate: stringPointer("2023-01-01"),
		Description: stringPointer("A test movie"),
		CoverURL:    stringPointer("https://example.com/cover.jpg"),
		Runtime:     int32Pointer(120),
		Budget:      int32Pointer(1000000),
		BoxOffice:   int32Pointer(5000000),
	}

	// Test creating a movie
	movie, err := repo.CreateMovie(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create movie: %v", err)
	}

	// Verify movie properties
	if movie.Title != input.Title {
		t.Errorf("Expected title %s, got %s", input.Title, movie.Title)
	}
	if movie.ReleaseDate == nil || *movie.ReleaseDate != *input.ReleaseDate {
		t.Errorf("Expected release date %v, got %v", input.ReleaseDate, movie.ReleaseDate)
	}
	if movie.Description == nil || *movie.Description != *input.Description {
		t.Errorf("Expected description %v, got %v", input.Description, movie.Description)
	}
	if movie.ID == uuid.Nil {
		t.Error("Expected non-nil movie ID")
	}
}

func TestMediaRepository_GetMovieByID(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create a test movie first
	input := model.CreateMovieInput{
		Title:       GenerateTestTitle("movie"),
		Description: stringPointer("Test movie for retrieval"),
		Runtime:     int32Pointer(90),
	}

	createdMovie, err := repo.CreateMovie(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	// Test getting movie by ID
	movie, err := repo.GetMovieByID(ctx, createdMovie.ID)
	if err != nil {
		t.Fatalf("Failed to get movie by ID: %v", err)
	}

	// Verify movie properties
	if movie.ID != createdMovie.ID {
		t.Errorf("Expected ID %v, got %v", createdMovie.ID, movie.ID)
	}
	if movie.Title != createdMovie.Title {
		t.Errorf("Expected title %s, got %s", createdMovie.Title, movie.Title)
	}

	// Test non-existent movie
	nonExistentID := uuid.New()
	_, err = repo.GetMovieByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent movie")
	}
}

func TestMediaRepository_GetAllMovies(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create multiple test movies
	movies := make([]*model.Movie, 3)
	for i := 0; i < 3; i++ {
		input := model.CreateMovieInput{
			Title:       GenerateTestTitle("movie"),
			Description: stringPointer("Test movie for listing"),
		}

		movie, err := repo.CreateMovie(ctx, input)
		if err != nil {
			t.Fatalf("Failed to create test movie %d: %v", i, err)
		}
		movies[i] = movie
	}

	// Test getting all movies
	allMovies, err := repo.GetAllMovies(ctx)
	if err != nil {
		t.Fatalf("Failed to get all movies: %v", err)
	}

	// Should have at least our test movies
	if len(allMovies) < 3 {
		t.Errorf("Expected at least 3 movies, got %d", len(allMovies))
	}

	// Verify our test movies are in the results
	foundMovies := 0
	for _, createdMovie := range movies {
		for _, fetchedMovie := range allMovies {
			if fetchedMovie.ID == createdMovie.ID {
				foundMovies++
				break
			}
		}
	}

	if foundMovies != 3 {
		t.Errorf("Expected to find all 3 created movies, found %d", foundMovies)
	}
}

func TestMediaRepository_CreateBook(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateBookInput{
		Title:       GenerateTestTitle("book"),
		ReleaseDate: stringPointer("2023-06-01"),
		Description: stringPointer("A test book"),
		CoverURL:    stringPointer("https://example.com/book-cover.jpg"),
		Pages:       int32Pointer(300),
		Isbn:        stringPointer("978-0123456789"),
		Publisher:   stringPointer("Test Publishers"),
	}

	// Test creating a book
	book, err := repo.CreateBook(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Verify book properties
	if book.Title != input.Title {
		t.Errorf("Expected title %s, got %s", input.Title, book.Title)
	}
	if book.Pages == nil || *book.Pages != *input.Pages {
		t.Errorf("Expected pages %v, got %v", input.Pages, book.Pages)
	}
	if book.Isbn == nil || *book.Isbn != *input.Isbn {
		t.Errorf("Expected ISBN %v, got %v", input.Isbn, book.Isbn)
	}
	if book.ID == uuid.Nil {
		t.Error("Expected non-nil book ID")
	}
}

func TestMediaRepository_CreateGame(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateGameInput{
		Title:       GenerateTestTitle("game"),
		ReleaseDate: stringPointer("2023-03-15"),
		Description: stringPointer("A test game"),
		CoverURL:    stringPointer("https://example.com/game-cover.jpg"),
		Genre:       []string{"Action", "Adventure"},
		EsrbRating:  stringPointer("T"),
		Multiplayer: boolPointer(true),
	}

	// Test creating a game
	game, err := repo.CreateGame(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Verify game properties
	if game.Title != input.Title {
		t.Errorf("Expected title %s, got %s", input.Title, game.Title)
	}
	if len(game.Genre) != len(input.Genre) {
		t.Errorf("Expected %d genres, got %d", len(input.Genre), len(game.Genre))
	}
	if game.EsrbRating == nil || *game.EsrbRating != *input.EsrbRating {
		t.Errorf("Expected ESRB rating %v, got %v", input.EsrbRating, game.EsrbRating)
	}
	if game.ID == uuid.Nil {
		t.Error("Expected non-nil game ID")
	}
}

func TestMediaRepository_CreateTVShow(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateTVShowInput{
		Title:       GenerateTestTitle("tvshow"),
		ReleaseDate: stringPointer("2023-09-01"),
		Description: stringPointer("A test TV show"),
		CoverURL:    stringPointer("https://example.com/tv-cover.jpg"),
		Seasons:     int32Pointer(3),
		Episodes:    int32Pointer(30),
		Status:      stringPointer("Ongoing"),
	}

	// Test creating a TV show
	tvShow, err := repo.CreateTVShow(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create TV show: %v", err)
	}

	// Verify TV show properties
	if tvShow.Title != input.Title {
		t.Errorf("Expected title %s, got %s", input.Title, tvShow.Title)
	}
	if tvShow.Seasons == nil || *tvShow.Seasons != *input.Seasons {
		t.Errorf("Expected seasons %v, got %v", input.Seasons, tvShow.Seasons)
	}
	if tvShow.Status == nil || *tvShow.Status != *input.Status {
		t.Errorf("Expected status %v, got %v", input.Status, tvShow.Status)
	}
	if tvShow.ID == uuid.Nil {
		t.Error("Expected non-nil TV show ID")
	}
}

func TestMediaRepository_CreateMusicAlbum(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Test data
	input := model.CreateMusicAlbumInput{
		Title:       GenerateTestTitle("album"),
		ReleaseDate: stringPointer("2023-12-01"),
		Description: stringPointer("A test music album"),
		CoverURL:    stringPointer("https://example.com/album-cover.jpg"),
		TrackCount:  int32Pointer(12),
		Duration:    int32Pointer(3600), // 60 minutes in seconds
		Label:       stringPointer("Test Records"),
	}

	// Test creating a music album
	album, err := repo.CreateMusicAlbum(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create music album: %v", err)
	}

	// Verify music album properties
	if album.Title != input.Title {
		t.Errorf("Expected title %s, got %s", input.Title, album.Title)
	}
	if album.TrackCount == nil || *album.TrackCount != *input.TrackCount {
		t.Errorf("Expected track count %v, got %v", input.TrackCount, album.TrackCount)
	}
	if album.Label == nil || *album.Label != *input.Label {
		t.Errorf("Expected label %v, got %v", input.Label, album.Label)
	}
	if album.ID == uuid.Nil {
		t.Error("Expected non-nil music album ID")
	}
}

func TestMediaRepository_GetMediaByID(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create different types of media
	movieInput := model.CreateMovieInput{
		Title: GenerateTestTitle("movie"),
	}
	movie, err := repo.CreateMovie(ctx, movieInput)
	if err != nil {
		t.Fatalf("Failed to create test movie: %v", err)
	}

	bookInput := model.CreateBookInput{
		Title: GenerateTestTitle("book"),
	}
	book, err := repo.CreateBook(ctx, bookInput)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	// Test getting movie by generic media ID
	media, err := repo.GetMediaByID(ctx, movie.ID)
	if err != nil {
		t.Fatalf("Failed to get movie as media: %v", err)
	}
	if media.GetID() != movie.ID {
		t.Errorf("Expected media ID %v, got %v", movie.ID, media.GetID())
	}

	// Test getting book by generic media ID
	media, err = repo.GetMediaByID(ctx, book.ID)
	if err != nil {
		t.Fatalf("Failed to get book as media: %v", err)
	}
	if media.GetID() != book.ID {
		t.Errorf("Expected media ID %v, got %v", book.ID, media.GetID())
	}

	// Test non-existent media
	nonExistentID := uuid.New()
	_, err = repo.GetMediaByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent media")
	}
}

func TestMediaRepository_GetAllMedia(t *testing.T) {
	repo, testDB := setupTestRepository(t)
	defer testDB.Close(t)

	ctx := context.Background()

	// Create different types of media
	createdMedia := make([]uuid.UUID, 0, 6)

	// Movies
	for i := 0; i < 2; i++ {
		input := model.CreateMovieInput{
			Title: GenerateTestTitle("movie"),
		}
		movie, err := repo.CreateMovie(ctx, input)
		if err != nil {
			t.Fatalf("Failed to create test movie %d: %v", i, err)
		}
		createdMedia = append(createdMedia, movie.ID)
	}

	// Books
	for i := 0; i < 2; i++ {
		input := model.CreateBookInput{
			Title: GenerateTestTitle("book"),
		}
		book, err := repo.CreateBook(ctx, input)
		if err != nil {
			t.Fatalf("Failed to create test book %d: %v", i, err)
		}
		createdMedia = append(createdMedia, book.ID)
	}

	// Games
	for i := 0; i < 2; i++ {
		input := model.CreateGameInput{
			Title: GenerateTestTitle("game"),
		}
		game, err := repo.CreateGame(ctx, input)
		if err != nil {
			t.Fatalf("Failed to create test game %d: %v", i, err)
		}
		createdMedia = append(createdMedia, game.ID)
	}

	// Test getting all media
	allMedia, err := repo.GetAllMedia(ctx)
	if err != nil {
		t.Fatalf("Failed to get all media: %v", err)
	}

	// Should have at least our test media
	if len(allMedia) < 6 {
		t.Errorf("Expected at least 6 media items, got %d", len(allMedia))
	}

	// Verify our test media are in the results
	foundMedia := 0
	for _, createdID := range createdMedia {
		for _, fetchedMedia := range allMedia {
			if fetchedMedia.GetID() == createdID {
				foundMedia++
				break
			}
		}
	}

	if foundMedia != 6 {
		t.Errorf("Expected to find all 6 created media items, found %d", foundMedia)
	}
}

// Helper functions for pointer types
func int32Pointer(i int32) *int32 {
	return &i
}

func float64Pointer(f float64) *float64 {
	return &f
}

func boolPointer(b bool) *bool {
	return &b
}
