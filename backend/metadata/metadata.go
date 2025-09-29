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

// NewService creates a new metadata service with the default fetchers
func NewService() (*Service, error) {
	s := &Service{
		fetchers: make(map[MediaType]Fetcher),
	}

	// Initialize fetchers for each media type
	// Note: You'll need to set the required environment variables for these to work
	if err := s.initFetchers(); err != nil {
		return nil, fmt.Errorf("failed to initialize fetchers: %w", err)
	}

	return s, nil
}

// initFetchers initializes all available fetchers
func (s *Service) initFetchers() error {
	// Initialize book fetcher
	bookFetcher, err := NewBookFetcher()
	if err != nil {
		return fmt.Errorf("failed to initialize book fetcher: %w", err)
	}
	s.fetchers[MediaTypeBook] = bookFetcher

	// Initialize game fetcher
	gameFetcher, err := NewGameFetcher()
	if err != nil {
		return fmt.Errorf("failed to initialize game fetcher: %w", err)
	}
	s.fetchers[MediaTypeGame] = gameFetcher

	// Initialize video fetcher (for both movies and TV shows)
	videoFetcher, err := NewVideoFetcher()
	if err != nil {
		return fmt.Errorf("failed to initialize video fetcher: %w", err)
	}
	s.fetchers[MediaTypeMovie] = videoFetcher
	s.fetchers[MediaTypeTV] = videoFetcher

	return nil
}

// GetMetadata fetches metadata for the given media info
func (s *Service) GetMetadata(info MediaInfo) (*MediaMetadata, error) {
	if info.Type == "" {
		return nil, errors.New("media type is required")
	}

	fetcher, exists := s.fetchers[info.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported media type: %s", info.Type)
	}

	return fetcher.Fetch(info)
}

// GetMetadata is a convenience function that creates a new service and fetches metadata
func GetMetadata(mediaType string, mediaInfo map[string]interface{}) (*MediaMetadata, error) {
	service, err := NewService()
	if err != nil {
		return nil, err
	}

	info := MediaInfo{
		Type:  MediaType(mediaType),
		Title: mediaInfo["title"].(string),
		ID:    mediaInfo["id"].(string),
	}

	if year, ok := mediaInfo["year"].(int); ok {
		info.Year = year
	}

	return service.GetMetadata(info)
}
