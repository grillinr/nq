// Package integrations handles third-party service integrations and authentication.
package integrations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
)

// Integration represents a third-party service integration
type Integration interface {
	// GetName returns the integration name
	GetName() string

	// Authenticate handles the authentication process
	Authenticate(ctx context.Context, credentials map[string]string) error

	// SyncUserData fetches and syncs user data from the service
	SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error)

	// IsAuthenticated checks if the integration is properly authenticated
	IsAuthenticated() bool

	// GetSupportedMediaTypes returns the media types this integration supports
	GetSupportedMediaTypes() []MediaType
}

// MediaType represents different types of media content
type MediaType string

const (
	MediaTypeGame    MediaType = "game"
	MediaTypeMusic   MediaType = "music"
	MediaTypeVideo   MediaType = "video"
	MediaTypeBook    MediaType = "book"
	MediaTypeArticle MediaType = "article"
)

// SyncResult contains the results of a data synchronization operation
type SyncResult struct {
	IntegrationName string                      `json:"integration_name"`
	UserID          uuid.UUID                   `json:"user_id"`
	SyncedAt        time.Time                   `json:"synced_at"`
	ItemsProcessed  int                         `json:"items_processed"`
	ItemsAdded      int                         `json:"items_added"`
	ItemsUpdated    int                         `json:"items_updated"`
	Errors          []string                    `json:"errors,omitempty"`
	MediaData       map[MediaType][]interface{} `json:"media_data"`
}

// Manager handles all third-party integrations
type Manager struct {
	integrations map[string]Integration
}

// NewManager creates a new integration manager
func NewManager() *Manager {
	return &Manager{
		integrations: make(map[string]Integration),
	}
}

// RegisterIntegration registers a new integration
func (m *Manager) RegisterIntegration(integration Integration) {
	m.integrations[integration.GetName()] = integration
}

// GetIntegration retrieves an integration by name
func (m *Manager) GetIntegration(name string) (Integration, error) {
	integration, exists := m.integrations[name]
	if !exists {
		return nil, fmt.Errorf("integration %s not found", name)
	}
	return integration, nil
}

// ListIntegrations returns all available integrations
func (m *Manager) ListIntegrations() []Integration {
	integrations := make([]Integration, 0, len(m.integrations))
	for _, integration := range m.integrations {
		integrations = append(integrations, integration)
	}
	return integrations
}

// SyncAllUserData syncs data from all authenticated integrations for a user
func (m *Manager) SyncAllUserData(ctx context.Context, userID uuid.UUID) (map[string]*SyncResult, error) {
	results := make(map[string]*SyncResult)
	var errors []string

	for name, integration := range m.integrations {
		if !integration.IsAuthenticated() {
			continue
		}

		result, err := integration.SyncUserData(ctx, userID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
			continue
		}

		results[name] = result
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("sync errors occurred: %v", errors)
	}

	return results, nil
}

// BaseIntegration provides common functionality for all integrations
type BaseIntegration struct {
	name           string
	authenticated  bool
	supportedTypes []MediaType
	credentials    map[string]string
}

// NewBaseIntegration creates a new base integration
func NewBaseIntegration(name string, supportedTypes []MediaType) *BaseIntegration {
	return &BaseIntegration{
		name:           name,
		authenticated:  false,
		supportedTypes: supportedTypes,
		credentials:    make(map[string]string),
	}
}

// GetName returns the integration name
func (b *BaseIntegration) GetName() string {
	return b.name
}

// IsAuthenticated checks if the integration is authenticated
func (b *BaseIntegration) IsAuthenticated() bool {
	return b.authenticated
}

// GetSupportedMediaTypes returns supported media types
func (b *BaseIntegration) GetSupportedMediaTypes() []MediaType {
	return b.supportedTypes
}

// SetAuthenticated sets the authentication status
func (b *BaseIntegration) SetAuthenticated(authenticated bool) {
	b.authenticated = authenticated
}

// SetCredential sets a credential value
func (b *BaseIntegration) SetCredential(key, value string) {
	b.credentials[key] = value
}

// GetCredential gets a credential value
func (b *BaseIntegration) GetCredential(key string) (string, bool) {
	value, exists := b.credentials[key]
	return value, exists
}

// ValidateCredentials validates that required credentials are present
func (b *BaseIntegration) ValidateCredentials(requiredKeys []string) error {
	for _, key := range requiredKeys {
		if _, exists := b.credentials[key]; !exists {
			return fmt.Errorf("missing required credential: %s", key)
		}
	}
	return nil
}

// ConvertToMediaItem is a helper to convert integration data to standardized media items
func ConvertToMediaItem(mediaType MediaType, data map[string]interface{}) (model.Media, error) {
	switch mediaType {
	case MediaTypeGame:
		return convertToGame(data)
	case MediaTypeMusic:
		return convertToMusicAlbum(data)
	case MediaTypeBook:
		return convertToBook(data)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}
}

// Helper functions for converting data to model types
func convertToGame(data map[string]interface{}) (*model.Game, error) {
	game := &model.Game{
		ID: uuid.New(),
	}

	if title, ok := data["title"].(string); ok {
		game.Title = title
	} else {
		return nil, errors.New("title is required for game")
	}

	if description, ok := data["description"].(string); ok {
		game.Description = &description
	}

	if coverURL, ok := data["cover_url"].(string); ok {
		game.CoverURL = &coverURL
	}

	if releaseDate, ok := data["release_date"].(string); ok {
		game.ReleaseDate = &releaseDate
	}

	if genres, ok := data["genres"].([]string); ok {
		game.Genre = genres
	}

	if esrbRating, ok := data["esrb_rating"].(string); ok {
		game.EsrbRating = &esrbRating
	}

	if multiplayer, ok := data["multiplayer"].(bool); ok {
		game.Multiplayer = &multiplayer
	}

	return game, nil
}

func convertToMusicAlbum(data map[string]interface{}) (*model.MusicAlbum, error) {
	album := &model.MusicAlbum{
		ID: uuid.New(),
	}

	if title, ok := data["title"].(string); ok {
		album.Title = title
	} else {
		return nil, errors.New("title is required for music album")
	}

	if description, ok := data["description"].(string); ok {
		album.Description = &description
	}

	if coverURL, ok := data["cover_url"].(string); ok {
		album.CoverURL = &coverURL
	}

	if releaseDate, ok := data["release_date"].(string); ok {
		album.ReleaseDate = &releaseDate
	}

	if trackCount, ok := data["track_count"].(int32); ok {
		album.TrackCount = &trackCount
	}

	if duration, ok := data["duration"].(int32); ok {
		album.Duration = &duration
	}

	if label, ok := data["label"].(string); ok {
		album.Label = &label
	}

	return album, nil
}

func convertToBook(data map[string]interface{}) (*model.Book, error) {
	book := &model.Book{
		ID: uuid.New(),
	}

	if title, ok := data["title"].(string); ok {
		book.Title = title
	} else {
		return nil, errors.New("title is required for book")
	}

	if description, ok := data["description"].(string); ok {
		book.Description = &description
	}

	if coverURL, ok := data["cover_url"].(string); ok {
		book.CoverURL = &coverURL
	}

	if releaseDate, ok := data["release_date"].(string); ok {
		book.ReleaseDate = &releaseDate
	}

	if pages, ok := data["pages"].(int32); ok {
		book.Pages = &pages
	}

	if isbn, ok := data["isbn"].(string); ok {
		book.Isbn = &isbn
	}

	if publisher, ok := data["publisher"].(string); ok {
		book.Publisher = &publisher
	}

	return book, nil
}
