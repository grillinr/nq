package db

import (
	"context"
	"fmt"
	"time"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Helper to safely convert a Neo4j datetime value to RFC3339 string
func neo4jDateTimeToString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case time.Time:
		return v.Format(time.RFC3339)
	case *time.Time:
		if v == nil {
			return ""
		}
		return v.Format(time.RFC3339)
	case neo4j.LocalDateTime:
		// Convert LocalDateTime to time.Time then format
		t := v.Time()
		return t.Format(time.RFC3339)
	case neo4j.Time:
		t := v.Time()
		return t.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

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
				RatedAt: neo4jDateTimeToString(record.AsMap()["ratedAt"]),
			}
			if u, err := r.GetUserByID(ctx, userID); err == nil {
				rating.User = u
			}
			if m, err := r.GetMediaByID(ctx, mediaID); err == nil {
				rating.Media = m
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
				RatedAt: neo4jDateTimeToString(record.AsMap()["ratedAt"]),
			}
			if u, err := r.GetUserByID(ctx, userID); err == nil {
				rating.User = u
			}
			if m, err := r.GetMediaByID(ctx, mediaID); err == nil {
				rating.Media = m
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
				RatedAt: neo4jDateTimeToString(record.AsMap()["ratedAt"]),
			}
			if userIDStr, ok := record.AsMap()["userId"].(string); ok {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						rating.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						rating.Media = m
					}
				}
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
				RatedAt: neo4jDateTimeToString(record.AsMap()["ratedAt"]),
			}
			if userIDStr, ok := record.AsMap()["userId"].(string); ok {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						rating.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						rating.Media = m
					}
				}
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
				RatedAt: neo4jDateTimeToString(record.AsMap()["ratedAt"]),
			}
			if u, err := r.GetUserByID(ctx, userID); err == nil {
				rating.User = u
			}
			if m, err := r.GetMediaByID(ctx, mediaID); err == nil {
				rating.Media = m
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

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, err
		}

		// If no nodes/relationships deleted, rating did not exist
		if summary.Counters().NodesDeleted() == 0 && summary.Counters().RelationshipsDeleted() == 0 {
			return nil, fmt.Errorf("rating not found")
		}

		return summary, nil
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

			// Accept multiple numeric types
			switch v := avgValue.(type) {
			case float64:
				return &v, nil
			case int32:
				f := float64(v)
				return &f, nil
			case int64:
				f := float64(v)
				return &f, nil
			}
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// result may already be *float64 or nil
	switch v := result.(type) {
	case *float64:
		return v, nil
	case float64:
		return &v, nil
	default:
		return nil, nil
	}
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
