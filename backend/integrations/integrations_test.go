package integrations

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestIntegrationManager(t *testing.T) {
	manager := NewManager()

	// Test registering integrations
	spotify := NewSpotifyIntegration()
	manager.RegisterIntegration(spotify)

	steam := NewSteamIntegration()
	manager.RegisterIntegration(steam)

	youtube := NewYouTubeIntegration()
	manager.RegisterIntegration(youtube)

	// Test listing integrations
	integrations := manager.ListIntegrations()
	if len(integrations) != 3 {
		t.Errorf("Expected 3 integrations, got %d", len(integrations))
	}

	// Test getting specific integration
	spotifyRetrieved, err := manager.GetIntegration("spotify")
	if err != nil {
		t.Errorf("Failed to get Spotify integration: %v", err)
	}

	if spotifyRetrieved.GetName() != "spotify" {
		t.Errorf("Expected Spotify integration, got %s", spotifyRetrieved.GetName())
	}

	// Test getting non-existent integration
	_, err = manager.GetIntegration("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent integration")
	}
}

func TestBaseIntegration(t *testing.T) {
	base := NewBaseIntegration("test", []MediaType{MediaTypeMusic, MediaTypeGame})

	// Test basic properties
	if base.GetName() != "test" {
		t.Errorf("Expected name 'test', got %s", base.GetName())
	}

	if base.IsAuthenticated() {
		t.Error("Expected integration to start unauthenticated")
	}

	supportedTypes := base.GetSupportedMediaTypes()
	if len(supportedTypes) != 2 {
		t.Errorf("Expected 2 supported media types, got %d", len(supportedTypes))
	}

	// Test credentials
	base.SetCredential("api_key", "test-key")
	key, exists := base.GetCredential("api_key")
	if !exists || key != "test-key" {
		t.Errorf("Expected credential 'test-key', got %s (exists: %v)", key, exists)
	}

	// Test authentication
	base.SetAuthenticated(true)
	if !base.IsAuthenticated() {
		t.Error("Expected integration to be authenticated")
	}

	// Test credential validation
	err := base.ValidateCredentials([]string{"api_key"})
	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}

	err = base.ValidateCredentials([]string{"missing_key"})
	if err == nil {
		t.Error("Expected validation to fail for missing credential")
	}
}

func TestSpotifyIntegration(t *testing.T) {
	spotify := NewSpotifyIntegration()

	// Test basic properties
	if spotify.GetName() != "spotify" {
		t.Errorf("Expected name 'spotify', got %s", spotify.GetName())
	}

	supportedTypes := spotify.GetSupportedMediaTypes()
	if len(supportedTypes) != 1 || supportedTypes[0] != MediaTypeMusic {
		t.Errorf("Expected Spotify to support only music, got %v", supportedTypes)
	}

	if spotify.IsAuthenticated() {
		t.Error("Expected Spotify to start unauthenticated")
	}
}

func TestSteamIntegration(t *testing.T) {
	steam := NewSteamIntegration()

	// Test basic properties
	if steam.GetName() != "steam" {
		t.Errorf("Expected name 'steam', got %s", steam.GetName())
	}

	supportedTypes := steam.GetSupportedMediaTypes()
	if len(supportedTypes) != 1 || supportedTypes[0] != MediaTypeGame {
		t.Errorf("Expected Steam to support only games, got %v", supportedTypes)
	}
}

func TestYouTubeIntegration(t *testing.T) {
	youtube := NewYouTubeIntegration()

	// Test basic properties
	if youtube.GetName() != "youtube" {
		t.Errorf("Expected name 'youtube', got %s", youtube.GetName())
	}

	supportedTypes := youtube.GetSupportedMediaTypes()
	if len(supportedTypes) != 1 || supportedTypes[0] != MediaTypeVideo {
		t.Errorf("Expected YouTube to support only video, got %v", supportedTypes)
	}
}

func TestConvertToMediaItem(t *testing.T) {
	// Test game conversion
	gameData := map[string]interface{}{
		"title":       "Test Game",
		"description": "A test game",
		"cover_url":   "https://example.com/cover.jpg",
	}

	mediaItem, err := ConvertToMediaItem(MediaTypeGame, gameData)
	if err != nil {
		t.Errorf("Failed to convert game data: %v", err)
	}

	if mediaItem.GetTitle() != "Test Game" {
		t.Errorf("Expected title 'Test Game', got %s", mediaItem.GetTitle())
	}

	// Test music conversion
	musicData := map[string]interface{}{
		"title":       "Test Album",
		"track_count": int32(12),
		"label":       "Test Records",
	}

	mediaItem, err = ConvertToMediaItem(MediaTypeMusic, musicData)
	if err != nil {
		t.Errorf("Failed to convert music data: %v", err)
	}

	if mediaItem.GetTitle() != "Test Album" {
		t.Errorf("Expected title 'Test Album', got %s", mediaItem.GetTitle())
	}

	// Test book conversion
	bookData := map[string]interface{}{
		"title":     "Test Book",
		"isbn":      "1234567890",
		"pages":     int32(200),
		"publisher": "Test Publishing",
	}

	mediaItem, err = ConvertToMediaItem(MediaTypeBook, bookData)
	if err != nil {
		t.Errorf("Failed to convert book data: %v", err)
	}

	if mediaItem.GetTitle() != "Test Book" {
		t.Errorf("Expected title 'Test Book', got %s", mediaItem.GetTitle())
	}

	// Test unsupported media type
	_, err = ConvertToMediaItem(MediaType("unsupported"), gameData)
	if err == nil {
		t.Error("Expected error for unsupported media type")
	}

	// Test missing title
	invalidData := map[string]interface{}{
		"description": "Missing title",
	}

	_, err = ConvertToMediaItem(MediaTypeGame, invalidData)
	if err == nil {
		t.Error("Expected error for missing title")
	}
}

func TestSyncAllUserData(t *testing.T) {
	manager := NewManager()
	ctx := context.Background()
	userID := uuid.New()

	// Add some mock integrations (they won't be authenticated, so should be skipped)
	spotify := NewSpotifyIntegration()
	steam := NewSteamIntegration()

	manager.RegisterIntegration(spotify)
	manager.RegisterIntegration(steam)

	// Try to sync - should return empty results since no integrations are authenticated
	results, err := manager.SyncAllUserData(ctx, userID)
	if err != nil {
		t.Errorf("Expected no error for unauthenticated integrations, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results from unauthenticated integrations, got %d", len(results))
	}
}
