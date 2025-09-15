package db

import (
	"context"
	"fmt"
	"nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateRating creates a new rating in the database
func (r *Neo4jRepository) CreateRating(ctx context.Context, userID, mediaID uuid.UUID, score float64) (*model.Rating, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userID})
			MATCH (m:Media {id: $mediaID})
			CREATE (r:Rating {
				userId: $userID,
				mediaId: $mediaID,
				score: $score,
				ratedAt: datetime()
			})
			CREATE (u)-[:RATED]->(r)
			CREATE (r)-[:RATING_FOR]->(m)
			RETURN r.userId as userId, r.mediaId as mediaId, r.score as score, r.ratedAt as ratedAt
		`

		params := map[string]any{
			"userID":  userID.String(),
			"mediaID": mediaID.String(),
			"score":   score,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			rating := &model.Rating{
				Score:   score,
				RatedAt: record.AsMap()["ratedAt"].(string),
				// TODO: Populate User and Media from IDs
			}
			return rating, nil
		}

		return nil, fmt.Errorf("failed to create rating")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Rating), nil
}

// GetRating retrieves a rating by user and media IDs
func (r *Neo4jRepository) GetRating(ctx context.Context, userID, mediaID uuid.UUID) (*model.Rating, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {userId: $userID, mediaId: $mediaID})
			RETURN r.userId as userId, r.mediaId as mediaId, r.score as score, r.ratedAt as ratedAt
		`

		params := map[string]any{
			"userID":  userID.String(),
			"mediaID": mediaID.String(),
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			rating := &model.Rating{
				Score:   getFloat64FromRecord(record, "score"),
				RatedAt: record.AsMap()["ratedAt"].(string),
				// TODO: Populate User and Media from IDs
			}
			return rating, nil
		}

		return nil, fmt.Errorf("rating not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Rating), nil
}

// GetUserRatings retrieves all ratings for a user
func (r *Neo4jRepository) GetUserRatings(ctx context.Context, userID uuid.UUID) ([]*model.Rating, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {userId: $userID})
			OPTIONAL MATCH (r)-[:RATING_FOR]->(m:Media)
			RETURN r.userId as userId, r.mediaId as mediaId, r.score as score, r.ratedAt as ratedAt
			ORDER BY r.ratedAt DESC
		`

		params := map[string]any{"userID": userID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var ratings []*model.Rating
		for result.Next(ctx) {
			record := result.Record()
			rating := &model.Rating{
				Score:   getFloat64FromRecord(record, "score"),
				RatedAt: record.AsMap()["ratedAt"].(string),
				// TODO: Populate User and Media from IDs
			}
			ratings = append(ratings, rating)
		}

		return ratings, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Rating), nil
}

// GetMediaRatings retrieves all ratings for a media item
func (r *Neo4jRepository) GetMediaRatings(ctx context.Context, mediaID uuid.UUID) ([]*model.Rating, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {mediaId: $mediaID})
			OPTIONAL MATCH (u:User)-[:RATED]->(r)
			RETURN r.userId as userId, r.mediaId as mediaId, r.score as score, r.ratedAt as ratedAt
			ORDER BY r.ratedAt DESC
		`

		params := map[string]any{"mediaID": mediaID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var ratings []*model.Rating
		for result.Next(ctx) {
			record := result.Record()
			rating := &model.Rating{
				Score:   getFloat64FromRecord(record, "score"),
				RatedAt: record.AsMap()["ratedAt"].(string),
				// TODO: Populate User and Media from IDs
			}
			ratings = append(ratings, rating)
		}

		return ratings, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Rating), nil
}

// UpdateRating updates an existing rating
func (r *Neo4jRepository) UpdateRating(ctx context.Context, userID, mediaID uuid.UUID, score float64) (*model.Rating, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {userId: $userID, mediaId: $mediaID})
			SET r.score = $score, r.ratedAt = datetime()
			RETURN r.userId as userId, r.mediaId as mediaId, r.score as score, r.ratedAt as ratedAt
		`

		params := map[string]any{
			"userID":  userID.String(),
			"mediaID": mediaID.String(),
			"score":   score,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			rating := &model.Rating{
				Score:   score,
				RatedAt: record.AsMap()["ratedAt"].(string),
				// TODO: Populate User and Media from IDs
			}
			return rating, nil
		}

		return nil, fmt.Errorf("rating not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Rating), nil
}

// DeleteRating deletes a rating
func (r *Neo4jRepository) DeleteRating(ctx context.Context, userID, mediaID uuid.UUID) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {userId: $userID, mediaId: $mediaID})
			DETACH DELETE r
		`

		params := map[string]any{
			"userID":  userID.String(),
			"mediaID": mediaID.String(),
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return result.Consume(ctx)
	})

	return err
}

// GetAverageRating calculates the average rating for a media item
func (r *Neo4jRepository) GetAverageRating(ctx context.Context, mediaID uuid.UUID) (*float64, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (r:Rating {mediaId: $mediaID})
			RETURN avg(r.score) as averageRating
		`

		params := map[string]any{"mediaID": mediaID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			avgValue := record.AsMap()["averageRating"]
			if avgValue == nil {
				return nil, nil
			}

			if avg, ok := avgValue.(float64); ok {
				return &avg, nil
			}
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*float64), nil
}

// Helper function to safely get float64 from record
func getFloat64FromRecord(record *neo4j.Record, key string) float64 {
	value := record.AsMap()[key]
	if value == nil {
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	}

	return 0.0
}
