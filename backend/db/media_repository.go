package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateMusicAlbum creates a new music album in the database
func (r *Neo4jRepository) CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error) {
	albumID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (a:MusicAlbum:Media {
				id: $id,
				title: $title,
				releaseDate: $releaseDate,
				description: $description,
				coverUrl: $coverUrl,
				trackCount: $trackCount,
				duration: $duration,
				label: $label
			})
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
		`

		params := map[string]any{
			"id":          albumID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"trackCount":  input.TrackCount,
			"duration":    input.Duration,
			"label":       input.Label,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			album := &model.MusicAlbum{
				ID:          albumID,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			return album, nil
		}

		if err := result.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("create music album returned no records")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.MusicAlbum), nil
}

// GetMusicAlbumByID retrieves a music album by its ID
func (r *Neo4jRepository) GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:MusicAlbum {id: $id})
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			album := &model.MusicAlbum{
				ID:          id,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			return album, nil
		}

		return nil, fmt.Errorf("music album not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.MusicAlbum), nil
}

// GetAllMusicAlbums retrieves all music albums
func (r *Neo4jRepository) GetAllMusicAlbums(ctx context.Context) ([]*model.MusicAlbum, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:MusicAlbum:Media)
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
			ORDER BY a.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var albums []*model.MusicAlbum
		for result.Next(ctx) {
			record := result.Record()
			idStr, _ := record.Get("id")
			albumID, _ := uuid.Parse(idStr.(string))
			album := &model.MusicAlbum{
				ID:          albumID,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			albums = append(albums, album)
		}

		return albums, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.MusicAlbum), nil
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
	if album, err := r.GetMusicAlbumByID(ctx, id); err == nil {
		return album, nil
	}

	return nil, fmt.Errorf("media not found")
}

// GetAllMedia retrieves all media items
func (r *Neo4jRepository) GetAllMedia(ctx context.Context) ([]model.Media, error) {
	var mediaItems []model.Media

	movies, err := r.GetAllMovies(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range movies {
		mediaItems = append(mediaItems, m)
	}

	books, err := r.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		mediaItems = append(mediaItems, b)
	}

	games, err := r.GetAllGames(ctx)
	if err != nil {
		return nil, err
	}
	for _, g := range games {
		mediaItems = append(mediaItems, g)
	}

	tvShows, err := r.GetAllTVShows(ctx)
	if err != nil {
		return nil, err
	}
	for _, tv := range tvShows {
		mediaItems = append(mediaItems, tv)
	}

	albums, err := r.GetAllMusicAlbums(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range albums {
		mediaItems = append(mediaItems, a)
	}

	return mediaItems, nil
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
