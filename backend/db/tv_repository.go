package db

import (
	"context"
	"fmt"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateTVShow creates a new TV show in the database
func (r *Neo4jRepository) CreateTVShow(ctx context.Context, input model.CreateTVShowInput) (*model.TVShow, error) {
	// Try to enrich with metadata if minimal data provided
	if r.metadata != nil && shouldEnrichTVShow(input) {
		enrichedInput, err := r.enrichTVShowInput(input)
		if err != nil {
			// Log error but continue with original input
			// In production, use proper logging
		} else {
			input = enrichedInput
		}
	}

	tvShowID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (t:TVShow:Media {
				id: $id,
				title: $title,
				releaseDate: $releaseDate,
				description: $description,
				coverUrl: $coverUrl,
				seasons: $seasons,
				episodes: $episodes,
				status: $status,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
			       t.description as description, t.coverUrl as coverUrl,
			       t.seasons as seasons, t.episodes as episodes, t.status as status
		`

		params := map[string]any{
			"id":          tvShowID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"seasons":     input.Seasons,
			"episodes":    input.Episodes,
			"status":      input.Status,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			tvShow := &model.TVShow{
				ID:            tvShowID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return tvShow, nil
		}

		return nil, fmt.Errorf("failed to create TV show")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.TVShow), nil
}

// GetTVShowByID retrieves a TV show by its ID
func (r *Neo4jRepository) GetTVShowByID(ctx context.Context, id uuid.UUID) (*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
				MATCH (t:TVShow {id: $id})
				OPTIONAL MATCH (t)<-[ract:ACTED_IN]-(cast:Person)
				OPTIONAL MATCH (t)<-[rcrew:CREW_ON]-(crew:Person)
				RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
				       t.description as description, t.coverUrl as coverUrl,
				       t.seasons as seasons, t.episodes as episodes, t.status as status,
				       collect(DISTINCT cast) as cast,
				       collect(DISTINCT crew) as crew,
				       collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
				       collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits
			`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			castCredits := parsePersonCredits(record.AsMap()["castCredits"])
			crewCredits := parseCrewCredits(record.AsMap()["crewCredits"])
			tvShow := &model.TVShow{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				CastCredits:   castCredits,
				CrewCredits:   crewCredits,
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return tvShow, nil
		}

		return nil, fmt.Errorf("TV show not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.TVShow), nil
}

// GetAllTVShows retrieves all TV shows
func (r *Neo4jRepository) GetAllTVShows(ctx context.Context) ([]*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (t:TVShow)
			RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
			       t.description as description, t.coverUrl as coverUrl,
			       t.seasons as seasons, t.episodes as episodes, t.status as status
			ORDER BY t.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var tvShows []*model.TVShow
		for result.Next(ctx) {
			record := result.Record()
			tvShowID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			tvShow := &model.TVShow{
				ID:            tvShowID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			tvShows = append(tvShows, tvShow)
		}

		return tvShows, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.TVShow), nil
}

// shouldEnrichTVShow determines if a TV show input should be enriched with metadata
func shouldEnrichTVShow(input model.CreateTVShowInput) bool {
	return input.Description == nil && input.CoverURL == nil && input.Seasons == nil
}

// enrichTVShowInput fetches metadata and merges it with the input
func (r *Neo4jRepository) enrichTVShowInput(input model.CreateTVShowInput) (model.CreateTVShowInput, error) {
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
		Type:        metadata.MediaTypeTV,
		Title:       input.Title,
		ReleaseYear: year,
	})

	if err != nil {
		return input, err
	}

	meta, ok := metaInterface.(*metadata.VideoMetadata)
	if !ok {
		return input, fmt.Errorf("unexpected metadata type for TV show")
	}

	// Merge metadata with input (input takes precedence)
	enriched := input

	if enriched.Description == nil && meta.Description != "" {
		enriched.Description = &meta.Description
	}

	if enriched.CoverURL == nil && meta.ImageURL != "" {
		enriched.CoverURL = &meta.ImageURL
	}

	if enriched.ReleaseDate == nil && meta.ReleaseYear > 0 {
		releaseDate := fmt.Sprintf("%d", meta.ReleaseYear)
		enriched.ReleaseDate = &releaseDate
	}

	// Note: TMDB doesn't provide seasons/episodes/status in basic metadata
	// Those would need to be fetched separately if required

	return enriched, nil
}
