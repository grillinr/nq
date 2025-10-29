package db

import (
	"context"
	"fmt"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateRecommendation creates a new recommendation in the database
func (r *Neo4jRepository) CreateRecommendation(ctx context.Context, userID, mediaID uuid.UUID, recommenderID *uuid.UUID, source *string, score *float64) (*model.Recommendation, error) {
	recommendationID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userID})
			MATCH (m:Media {id: $mediaID})
			CREATE (rec:Recommendation {
				id: $recommendationID,
				userId: $userID,
				mediaId: $mediaID,
				recommenderId: $recommenderID,
				source: $source,
				score: $score,
				createdAt: datetime()
			})
			CREATE (u)-[:RECEIVED_RECOMMENDATION]->(rec)
			CREATE (rec)-[:RECOMMENDS]->(m)
		`

		params := map[string]any{
			"recommendationID": recommendationID.String(),
			"userID":           userID.String(),
			"mediaID":          mediaID.String(),
		}

		if recommenderID != nil {
			params["recommenderID"] = recommenderID.String()
		} else {
			params["recommenderID"] = nil
		}

		params["source"] = source
		params["score"] = score

		query += `
			RETURN rec.id as id, rec.userId as userId, rec.mediaId as mediaId,
				   rec.recommenderId as recommenderId, rec.source as source, rec.score as score
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			recommendation := &model.Recommendation{
				ID:     recommendationID,
				Source: getStringPointer(record.AsMap()["source"]),
				Score:  getFloat64Pointer(record.AsMap()["score"]),
			}
			// Populate nested fields where possible (non-fatal)
			if u, err := r.GetUserByID(ctx, userID); err == nil {
				recommendation.User = u
			}
			if m, err := r.GetMediaByID(ctx, mediaID); err == nil {
				recommendation.Media = m
			}
			if recommenderID != nil {
				if recUser, err := r.GetUserByID(ctx, *recommenderID); err == nil {
					recommendation.Recommender = recUser
				}
			}
			return recommendation, nil
		}

		return nil, fmt.Errorf("failed to create recommendation")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Recommendation), nil
}

// GetRecommendations retrieves all recommendations for a user
func (r *Neo4jRepository) GetRecommendations(ctx context.Context, userID uuid.UUID) ([]*model.Recommendation, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (rec:Recommendation {userId: $userID})
			OPTIONAL MATCH (rec)-[:RECOMMENDS]->(m:Media)
			OPTIONAL MATCH (rec)-[:RECOMMENDED_BY]->(ru:User)
			RETURN rec.id as id, rec.userId as userId, rec.mediaId as mediaId,
				   rec.recommenderId as recommenderId, rec.source as source, rec.score as score
			ORDER BY rec.score DESC, rec.createdAt DESC
		`

		params := map[string]any{"userID": userID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var recommendations []*model.Recommendation
		for result.Next(ctx) {
			record := result.Record()
			recommendationID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			recommendation := &model.Recommendation{
				ID:     recommendationID,
				Source: getStringPointer(record.AsMap()["source"]),
				Score:  getFloat64Pointer(record.AsMap()["score"]),
			}
			// attempt to populate nested fields using returned ids from record
			if userIDStr, ok := record.AsMap()["userId"].(string); ok {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						recommendation.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						recommendation.Media = m
					}
				}
			}
			if recommenderIDStr, ok := record.AsMap()["recommenderId"].(string); ok && recommenderIDStr != "" {
				if rid, err := uuid.Parse(recommenderIDStr); err == nil {
					if ru, err := r.GetUserByID(ctx, rid); err == nil {
						recommendation.Recommender = ru
					}
				}
			}
			recommendations = append(recommendations, recommendation)
		}

		return recommendations, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Recommendation), nil
}

// GetRecommendationByID retrieves a recommendation by its ID
func (r *Neo4jRepository) GetRecommendationByID(ctx context.Context, id uuid.UUID) (*model.Recommendation, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (rec:Recommendation {id: $id})
			OPTIONAL MATCH (rec)-[:RECOMMENDS]->(m:Media)
			OPTIONAL MATCH (rec)-[:RECOMMENDED_BY]->(ru:User)
			RETURN rec.id as id, rec.userId as userId, rec.mediaId as mediaId,
				   rec.recommenderId as recommenderId, rec.source as source, rec.score as score
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			recommendation := &model.Recommendation{
				ID:     id,
				Source: getStringPointer(record.AsMap()["source"]),
				Score:  getFloat64Pointer(record.AsMap()["score"]),
			}
			if userIDStr, ok := record.AsMap()["userId"].(string); ok {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						recommendation.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						recommendation.Media = m
					}
				}
			}
			if recommenderIDStr, ok := record.AsMap()["recommenderId"].(string); ok && recommenderIDStr != "" {
				if rid, err := uuid.Parse(recommenderIDStr); err == nil {
					if ru, err := r.GetUserByID(ctx, rid); err == nil {
						recommendation.Recommender = ru
					}
				}
			}
			return recommendation, nil
		}

		return nil, fmt.Errorf("recommendation not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Recommendation), nil
}

// DeleteRecommendation deletes a recommendation
func (r *Neo4jRepository) DeleteRecommendation(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (rec:Recommendation {id: $id})
			WITH rec, rec.id as id
			DETACH DELETE rec
			RETURN id
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		// If no rows returned then the recommendation did not exist
		if !result.Next(ctx) {
			return nil, fmt.Errorf("recommendation not found")
		}

		// consume remaining result to finish the query
		if _, err := result.Consume(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
