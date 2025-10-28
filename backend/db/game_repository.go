package db

import (
	"context"
	"fmt"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateGame creates a new game in the database
func (r *Neo4jRepository) CreateGame(ctx context.Context, input model.CreateGameInput) (*model.Game, error) {
	// Try to enrich with metadata if minimal data provided
	if r.metadata != nil && shouldEnrichGame(input) {
		enrichedInput, err := r.enrichGameInput(input)
		if err != nil {
			// Log error but continue with original input
			// In production, use proper logging
		} else {
			input = enrichedInput
		}
	}

	gameID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (g:Game:Media {
				id: $id,
				title: $title,
				releaseDate: $releaseDate,
				description: $description,
				coverUrl: $coverUrl,
				genre: $genre,
				esrbRating: $esrbRating,
				multiplayer: $multiplayer,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN g.id as id, g.title as title, g.releaseDate as releaseDate,
			       g.description as description, g.coverUrl as coverUrl,
			       g.genre as genre, g.esrbRating as esrbRating, g.multiplayer as multiplayer
		`

		params := map[string]any{
			"id":          gameID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"genre":       input.Genre,
			"esrbRating":  input.EsrbRating,
			"multiplayer": input.Multiplayer,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			game := &model.Game{
				ID:            gameID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Genre:         getStringSlice(record.AsMap()["genre"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return game, nil
		}

		return nil, fmt.Errorf("failed to create game")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Game), nil
}

// GetGameByID retrieves a game by its ID
func (r *Neo4jRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (g:Game {id: $id})
			RETURN g.id as id, g.title as title, g.releaseDate as releaseDate,
			       g.description as description, g.coverUrl as coverUrl,
			       g.genre as genre, g.esrbRating as esrbRating, g.multiplayer as multiplayer
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			game := &model.Game{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Genre:         getStringSlice(record.AsMap()["genre"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return game, nil
		}

		return nil, fmt.Errorf("game not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Game), nil
}

// GetAllGames retrieves all games
func (r *Neo4jRepository) GetAllGames(ctx context.Context) ([]*model.Game, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (g:Game)
			RETURN g.id as id, g.title as title, g.releaseDate as releaseDate,
			       g.description as description, g.coverUrl as coverUrl,
			       g.genre as genre, g.esrbRating as esrbRating, g.multiplayer as multiplayer
			ORDER BY g.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var games []*model.Game
		for result.Next(ctx) {
			record := result.Record()
			gameID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			game := &model.Game{
				ID:            gameID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Genre:         getStringSlice(record.AsMap()["genre"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			games = append(games, game)
		}

		return games, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.Game), nil
}

// shouldEnrichGame determines if a game input should be enriched with metadata
func shouldEnrichGame(input model.CreateGameInput) bool {
	return input.Description == nil && input.CoverURL == nil && len(input.Genre) == 0
}

// enrichGameInput fetches metadata and merges it with the input
func (r *Neo4jRepository) enrichGameInput(input model.CreateGameInput) (model.CreateGameInput, error) {
	if r.metadata == nil {
		return input, nil
	}

	// Determine year from release date if available
	var year int
	if input.ReleaseDate != nil {
		// Simple year extraction - in production, use proper date parsing
		if len(*input.ReleaseDate) >= 4 {
			fmt.Sscanf(*input.ReleaseDate, "%d", &year)
		}
	}

	metaInterface, err := r.metadata.GetMetadata(metadata.MediaInfo{
		Type:        metadata.MediaTypeGame,
		Title:       input.Title,
		ReleaseYear: year,
	})

	if err != nil {
		return input, err
	}

	meta, ok := metaInterface.(*metadata.MediaMetadata)
	if !ok {
		return input, fmt.Errorf("unexpected metadata type for game")
	}

	// Merge metadata with input (input takes precedence)
	enriched := input

	if enriched.Description == nil && meta.Description != "" {
		enriched.Description = &meta.Description
	}

	if enriched.CoverURL == nil && meta.ImageURL != "" {
		enriched.CoverURL = &meta.ImageURL
	}

	// Merge genres if input doesn't have any
	if len(enriched.Genre) == 0 && len(meta.Genres) > 0 {
		enriched.Genre = meta.Genres
	}

	// Note: IGDB doesn't provide ESRB rating or multiplayer info in basic metadata
	// Those would need to be fetched separately if required

	return enriched, nil
}

// Helper function to safely get string slice from interface{}
func getStringSlice(value interface{}) []string {
	if value == nil {
		return []string{}
	}

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			}
		}
		return result
	}

	return []string{}
}

// Helper function to safely get bool pointer from interface{}
func getBoolPointer(value interface{}) *bool {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case bool:
		return &v
	case string:
		if v == "true" {
			b := true
			return &b
		} else if v == "false" {
			b := false
			return &b
		}
	}

	return nil
}
