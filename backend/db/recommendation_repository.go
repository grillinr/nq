package db

import (
	"context"
	"fmt"
	"nq/graph/model"

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
			query += `
				MATCH (rec)-[:RECOMMENDED_BY]->(r:User {id: $recommenderID})
			`
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
				// TODO: Populate User, Media, and Recommender from IDs
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
			OPTIONAL MATCH (rec)-[:RECOMMENDED_BY]->(r:User)
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
				// TODO: Populate User, Media, and Recommender from IDs
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
			OPTIONAL MATCH (rec)-[:RECOMMENDED_BY]->(r:User)
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
				// TODO: Populate User, Media, and Recommender from IDs
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
			DETACH DELETE rec
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return result.Consume(ctx)
	})

	return err
}
