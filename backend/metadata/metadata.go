// Package metadata provides functions to retrieve metadata for various media types
// such as videos, books, and games using external APIs.
package metadata

import (
	"errors"
	"fmt"
)

// Service handles metadata fetching for different media types
type Service struct {
	fetchers map[MediaType]Fetcher
}

// IsValidMediaType checks if the given media type is supported
func IsValidMediaType(mediaType string) bool {
	switch MediaType(mediaType) {
	case MediaTypeBook, MediaTypeGame, MediaTypeMovie, MediaTypeTV:
		return true
	default:
		return false
	}
}

// NewService creates a new metadata service with the default fetchers
func NewService() (*Service, error) {
	s := &Service{
		fetchers: make(map[MediaType]Fetcher),
	}

	// Initialize fetchers for each media type
	if err := s.initFetchers(); err != nil {
		return nil, fmt.Errorf("failed to initialize fetchers: %w", err)
	}

	return s, nil
}

// initFetchers initializes all available fetchers
// Only initializes fetchers that have their required environment variables set
func (s *Service) initFetchers() error {
	var initErrors []string

	// Initialize book fetcher (no API key required)
	if bookFetcher, err := NewBookFetcher(); err != nil {
		initErrors = append(initErrors, fmt.Sprintf("book fetcher: %v", err))
	} else {
		s.fetchers[MediaTypeBook] = bookFetcher
	}

	// Initialize game fetcher (requires IGDB_CLIENT_ID)
	if gameFetcher, err := NewGameFetcher(); err != nil {
		// Don't treat missing API keys as fatal errors, just log them
		initErrors = append(initErrors, fmt.Sprintf("game fetcher: %v (games will not be available)", err))
	} else {
		s.fetchers[MediaTypeGame] = gameFetcher
	}

	// Initialize video fetcher (for both movies and TV shows)
	if videoFetcher, err := NewVideoFetcher(); err != nil {
		initErrors = append(initErrors, fmt.Sprintf("video fetcher: %v (movies/TV shows will not be available)", err))
	} else {
		s.fetchers[MediaTypeMovie] = videoFetcher
		s.fetchers[MediaTypeTV] = videoFetcher
	}

	// Only return error if no fetchers were initialized at all
	if len(s.fetchers) == 0 {
		return fmt.Errorf("no fetchers could be initialized: %v", initErrors)
	}

	return nil
}

// GetMetadata fetches metadata for the given media info
func (s *Service) GetMetadata(info MediaInfo) (*MediaMetadata, error) {
	if info.Type == "" {
		return nil, errors.New("media type is required")
	}

	if info.Title == "" && info.ID == "" {
		return nil, errors.New("either title or ID is required")
	}

	fetcher, exists := s.fetchers[info.Type]
	if !exists {
		if IsValidMediaType(string(info.Type)) {
			return nil, fmt.Errorf("fetcher for media type '%s' is not available (missing API credentials)", info.Type)
		}
		return nil, fmt.Errorf("unsupported media type: %s", info.Type)
	}

	return fetcher.Fetch(info)
}

// GetMetadata is a convenience function that creates a new service and fetches metadata
func GetMetadata(mediaType string, mediaInfo map[string]any) (*MediaMetadata, error) {
	service, err := NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata service: %w", err)
	}

	info := MediaInfo{
		Type: MediaType(mediaType),
	}

	// Safely extract title
	if title, ok := mediaInfo["title"]; ok {
		if titleStr, ok := title.(string); ok {
			info.Title = titleStr
		} else {
			return nil, errors.New("title must be a string")
		}
	}

	// Safely extract ID (optional)
	if id, ok := mediaInfo["id"]; ok {
		if idStr, ok := id.(string); ok {
			info.ID = idStr
		} else {
			return nil, errors.New("id must be a string")
		}
	}

	// Safely extract year (optional)
	if year, ok := mediaInfo["year"]; ok {
		switch v := year.(type) {
		case int:
			info.Year = v
		case float64:
			info.Year = int(v)
		case string:
			// Try to parse string as int
			if yearInt, parseErr := fmt.Sscanf(v, "%d", &info.Year); parseErr != nil || yearInt != 1 {
				return nil, errors.New("year must be a valid integer")
			}
		default:
			return nil, errors.New("year must be an integer")
		}
	}

	return service.GetMetadata(info)
}

// GetMetadataByTitle is a simpler convenience function for fetching metadata by title and type
func GetMetadataByTitle(mediaType, title string, year int) (*MediaMetadata, error) {
	if !IsValidMediaType(mediaType) {
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	if title == "" {
		return nil, errors.New("title is required")
	}

	service, err := NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata service: %w", err)
	}

	info := MediaInfo{
		Type:  MediaType(mediaType),
		Title: title,
		Year:  year,
	}

	return service.GetMetadata(info)
}

// GetMetadataByID is a convenience function for fetching metadata by ID and type
func GetMetadataByID(mediaType, id string) (*MediaMetadata, error) {
	if !IsValidMediaType(mediaType) {
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	if id == "" {
		return nil, errors.New("id is required")
	}

	service, err := NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata service: %w", err)
	}

	info := MediaInfo{
		Type: MediaType(mediaType),
		ID:   id,
	}

	return service.GetMetadata(info)
}
