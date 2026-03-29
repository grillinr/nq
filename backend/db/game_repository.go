package db

import (
	"context"
	"errors"
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

	// Generate UUID upfront - will be used if creating new, ignored if matching existing
	gameID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Handle searchDepth
		searchDepth := int32(0)
		if input.SearchDepth != nil {
			searchDepth = *input.SearchDepth
		}

		// Use MERGE to atomically check-and-create the game node, preventing race conditions
		// Note: Neo4j doesn't allow null in MERGE keys, so we use empty string for null dates
		releaseDate := ""
		if input.ReleaseDate != nil {
			releaseDate = *input.ReleaseDate
		}

		query := `
		MERGE (g:Game:Media {title: $title, releaseDate: $releaseDate})
		ON CREATE SET 
			g.id = $id,
			g.description = $description,
			g.coverUrl = $coverUrl,
			g.genre = $genre,
			g.themes = $themes,
			g.keywords = $keywords,
			g.gameModes = $gameModes,
			g.perspectives = $perspectives,
			g.franchises = $franchises,
			g.platformsList = $platforms,
			g.esrbRating = $esrbRating,
			g.multiplayer = $multiplayer,
			g.searchDepth = $searchDepth,
			g.createdAt = datetime(),
			g.updatedAt = datetime()
		ON MATCH SET
			g.searchDepth = CASE 
				WHEN $searchDepth < g.searchDepth THEN $searchDepth 
				ELSE g.searchDepth 
			END,
			g.updatedAt = datetime(),
			g.description = CASE WHEN g.description IS NULL AND $description IS NOT NULL THEN $description ELSE g.description END,
			g.coverUrl = CASE WHEN g.coverUrl IS NULL AND $coverUrl IS NOT NULL THEN $coverUrl ELSE g.coverUrl END,
			g.genre = CASE WHEN g.genre IS NULL AND $genre IS NOT NULL THEN $genre ELSE g.genre END,
			g.themes = CASE WHEN g.themes IS NULL AND $themes IS NOT NULL THEN $themes ELSE g.themes END,
			g.keywords = CASE WHEN g.keywords IS NULL AND $keywords IS NOT NULL THEN $keywords ELSE g.keywords END,
			g.gameModes = CASE WHEN g.gameModes IS NULL AND $gameModes IS NOT NULL THEN $gameModes ELSE g.gameModes END,
			g.perspectives = CASE WHEN g.perspectives IS NULL AND $perspectives IS NOT NULL THEN $perspectives ELSE g.perspectives END,
			g.franchises = CASE WHEN g.franchises IS NULL AND $franchises IS NOT NULL THEN $franchises ELSE g.franchises END,
			g.platformsList = CASE WHEN g.platformsList IS NULL AND $platforms IS NOT NULL THEN $platforms ELSE g.platformsList END,
			g.esrbRating = CASE WHEN g.esrbRating IS NULL AND $esrbRating IS NOT NULL THEN $esrbRating ELSE g.esrbRating END,
			g.multiplayer = CASE WHEN g.multiplayer IS NULL AND $multiplayer IS NOT NULL THEN $multiplayer ELSE g.multiplayer END
		WITH g
		FOREACH (tagName IN CASE WHEN $genre IS NULL THEN [] ELSE $genre END |
			MERGE (t:Tag {type: 'genre', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $themes IS NULL THEN [] ELSE $themes END |
			MERGE (t:Tag {type: 'theme', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $keywords IS NULL THEN [] ELSE $keywords END |
			MERGE (t:Tag {type: 'keyword', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $gameModes IS NULL THEN [] ELSE $gameModes END |
			MERGE (t:Tag {type: 'mode', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $perspectives IS NULL THEN [] ELSE $perspectives END |
			MERGE (t:Tag {type: 'perspective', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $franchises IS NULL THEN [] ELSE $franchises END |
			MERGE (t:Tag {type: 'franchise', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		FOREACH (tagName IN CASE WHEN $platforms IS NULL THEN [] ELSE $platforms END |
			MERGE (t:Tag {type: 'platform', normalizedName: toLower(trim(tagName))})
			ON CREATE SET t.id = randomUUID(), t.name = tagName
			MERGE (t)-[:TAGGED]->(g)
		)
		RETURN g.id as id, g.title as title, g.releaseDate as releaseDate,
		       g.description as description, g.coverUrl as coverUrl,
		       g.genre as genre, g.themes as themes, g.keywords as keywords,
		       g.gameModes as gameModes, g.perspectives as perspectives,
		       g.franchises as franchises, g.platformsList as platformsList,
		       g.esrbRating as esrbRating, g.multiplayer as multiplayer,
		       g.searchDepth as searchDepth
		`

		params := map[string]any{
			"id":           gameID.String(),
			"title":        input.Title,
			"releaseDate":  releaseDate,
			"description":  input.Description,
			"coverUrl":     input.CoverURL,
			"genre":        input.Genre,
			"themes":       input.Themes,
			"keywords":     input.Keywords,
			"gameModes":    input.GameModes,
			"perspectives": input.Perspectives,
			"franchises":   input.Franchises,
			"platforms":    input.Platforms,
			"esrbRating":   input.EsrbRating,
			"multiplayer":  input.Multiplayer,
			"searchDepth":  searchDepth,
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
				Themes:        getStringSlice(record.AsMap()["themes"]),
				Keywords:      getStringSlice(record.AsMap()["keywords"]),
				GameModes:     getStringSlice(record.AsMap()["gameModes"]),
				Perspectives:  getStringSlice(record.AsMap()["perspectives"]),
				Franchises:    getStringSlice(record.AsMap()["franchises"]),
				PlatformsList: getStringSlice(record.AsMap()["platformsList"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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
			       g.genre as genre, g.themes as themes, g.keywords as keywords,
			       g.gameModes as gameModes, g.perspectives as perspectives,
			       g.franchises as franchises, g.platformsList as platformsList,
			       g.esrbRating as esrbRating, g.multiplayer as multiplayer,
			       g.searchDepth as searchDepth
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
				Themes:        getStringSlice(record.AsMap()["themes"]),
				Keywords:      getStringSlice(record.AsMap()["keywords"]),
				GameModes:     getStringSlice(record.AsMap()["gameModes"]),
				Perspectives:  getStringSlice(record.AsMap()["perspectives"]),
				Franchises:    getStringSlice(record.AsMap()["franchises"]),
				PlatformsList: getStringSlice(record.AsMap()["platformsList"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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

// GetAllGames retrieves all games with optional pagination.
func (r *Neo4jRepository) GetAllGames(ctx context.Context, limit, offset *int) ([]*model.Game, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Paginate on bare Game nodes first, then project fields — consistent with
		// the pattern used by GetAllMovies, GetAllTVShows, and GetAllBooks.
		query := `
			MATCH (g:Game)
			WITH g ORDER BY g.title
		`

		params := map[string]any{}
		if offset != nil {
			query += " SKIP $offset"
			params["offset"] = *offset
		}
		if limit != nil {
			query += " LIMIT $limit"
			params["limit"] = *limit
		}

		query += `
			RETURN g.id as id, g.title as title, g.releaseDate as releaseDate,
			       g.description as description, g.coverUrl as coverUrl,
			       g.genre as genre, g.themes as themes, g.keywords as keywords,
			       g.gameModes as gameModes, g.perspectives as perspectives,
			       g.franchises as franchises, g.platformsList as platformsList,
			       g.esrbRating as esrbRating, g.multiplayer as multiplayer,
			       g.searchDepth as searchDepth
		`

		result, err := tx.Run(ctx, query, params)
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
				Themes:        getStringSlice(record.AsMap()["themes"]),
				Keywords:      getStringSlice(record.AsMap()["keywords"]),
				GameModes:     getStringSlice(record.AsMap()["gameModes"]),
				Perspectives:  getStringSlice(record.AsMap()["perspectives"]),
				Franchises:    getStringSlice(record.AsMap()["franchises"]),
				PlatformsList: getStringSlice(record.AsMap()["platformsList"]),
				EsrbRating:    getStringPointer(record.AsMap()["esrbRating"]),
				Multiplayer:   getBoolPointer(record.AsMap()["multiplayer"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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
	return input.Description == nil && input.CoverURL == nil && len(input.Genre) == 0 && len(input.Themes) == 0 && len(input.Keywords) == 0
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
		ID:          safeStringValue(input.ExternalID),
	})

	if err != nil {
		if errors.Is(err, metadata.ErrIGDBAuthFailed) {
			return input, nil
		}
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
	if len(enriched.Themes) == 0 && len(meta.Themes) > 0 {
		enriched.Themes = meta.Themes
	}
	if len(enriched.Keywords) == 0 && len(meta.Keywords) > 0 {
		enriched.Keywords = meta.Keywords
	}
	if len(enriched.GameModes) == 0 && len(meta.GameModes) > 0 {
		enriched.GameModes = meta.GameModes
	}
	if len(enriched.Perspectives) == 0 && len(meta.Perspectives) > 0 {
		enriched.Perspectives = meta.Perspectives
	}
	if len(enriched.Franchises) == 0 && len(meta.Franchises) > 0 {
		enriched.Franchises = meta.Franchises
	}
	if len(enriched.Platforms) == 0 && len(meta.Platforms) > 0 {
		enriched.Platforms = meta.Platforms
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
