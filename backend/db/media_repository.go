package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
)

// CreateMusicAlbum creates a new music album in the database
func (r *Neo4jRepository) CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error) {
	// Implementation similar to CreateMovie but for MusicAlbum
	return nil, fmt.Errorf("not implemented")
}

// GetMusicAlbumByID retrieves a music album by its ID
func (r *Neo4jRepository) GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetAllMusicAlbums retrieves all music albums
func (r *Neo4jRepository) GetAllMusicAlbums(ctx context.Context) ([]*model.MusicAlbum, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetMediaByID retrieves any media by its ID
func (r *Neo4jRepository) GetMediaByID(ctx context.Context, id uuid.UUID) (model.Media, error) {
	// Try each media type
	if movie, err := r.GetMovieByID(ctx, id); err == nil {
		return movie, nil
	}
	if tvShow, err := r.GetTVShowByID(ctx, id); err == nil {
		return tvShow, nil
	}
	if book, err := r.GetBookByID(ctx, id); err == nil {
		return book, nil
	}
	if game, err := r.GetGameByID(ctx, id); err == nil {
		return game, nil
	}
	// Music albums not implemented yet

	return nil, fmt.Errorf("media not found")
}

// GetAllMedia retrieves all media items
func (r *Neo4jRepository) GetAllMedia(ctx context.Context) ([]model.Media, error) {
	// Implementation to get all media types
	return nil, fmt.Errorf("not implemented")
}

// Helper function to safely get int32 pointer from interface{}
func getInt32Pointer(value interface{}) *int32 {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int32:
		return &v
	case int64:
		i32 := int32(v)
		return &i32
	case float64:
		i32 := int32(v)
		return &i32
	case string:
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			i32 := int32(i)
			return &i32
		}
	}

	return nil
}
