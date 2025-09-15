package db

import (
	"context"
	"fmt"
	"nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateActivity creates a new user activity in the database
func (r *Neo4jRepository) CreateActivity(ctx context.Context, input model.CreateActivityInput) (*model.UserActivity, error) {
	activityID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userID})
			MATCH (m:Media {id: $mediaID})
			CREATE (a:UserActivity {
				id: $activityID,
				statusId: $statusID,
				rating: $rating,
				review: $review,
				startedAt: $startedAt,
				finishedAt: $finishedAt,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			CREATE (u)-[:HAS_ACTIVITY]->(a)
			CREATE (a)-[:ACTIVITY_FOR]->(m)
			RETURN a.id as id, a.statusId as statusId, a.rating as rating,
			       a.review as review, a.startedAt as startedAt, a.finishedAt as finishedAt
		`

		params := map[string]any{
			"activityID": activityID.String(),
			"userID":     input.UserID.String(),
			"mediaID":    input.MediaID.String(),
			"statusID":   input.StatusID,
			"rating":     input.Rating,
			"review":     input.Review,
			"startedAt":  input.StartedAt,
			"finishedAt": input.FinishedAt,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			activity := &model.UserActivity{
				ID:         activityID,
				Status:     &model.ActivityStatus{ID: input.StatusID}, // TODO: Get actual status
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
				// TODO: Populate User and Media from IDs
			}
			return activity, nil
		}

		return nil, fmt.Errorf("failed to create activity")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.UserActivity), nil
}

// GetActivityByID retrieves an activity by its ID
func (r *Neo4jRepository) GetActivityByID(ctx context.Context, id uuid.UUID) (*model.UserActivity, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:UserActivity {id: $id})
			OPTIONAL MATCH (u:User)-[:HAS_ACTIVITY]->(a)
			OPTIONAL MATCH (a)-[:ACTIVITY_FOR]->(m:Media)
			RETURN a.id as id, a.statusId as statusId, a.rating as rating,
			       a.review as review, a.startedAt as startedAt, a.finishedAt as finishedAt,
			       u.id as userId, m.id as mediaId
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			activity := &model.UserActivity{
				ID:         id,
				Status:     &model.ActivityStatus{ID: getInt32FromRecord(record, "statusId")},
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
				// TODO: Populate User and Media
			}
			return activity, nil
		}

		return nil, fmt.Errorf("activity not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.UserActivity), nil
}

// GetUserActivities retrieves all activities for a user
func (r *Neo4jRepository) GetUserActivities(ctx context.Context, userID uuid.UUID) ([]*model.UserActivity, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userID})-[:HAS_ACTIVITY]->(a:UserActivity)
			OPTIONAL MATCH (a)-[:ACTIVITY_FOR]->(m:Media)
			RETURN a.id as id, a.statusId as statusId, a.rating as rating,
			       a.review as review, a.startedAt as startedAt, a.finishedAt as finishedAt,
			       m.id as mediaId
			ORDER BY a.createdAt DESC
		`

		params := map[string]any{"userID": userID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var activities []*model.UserActivity
		for result.Next(ctx) {
			record := result.Record()
			activityID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			activity := &model.UserActivity{
				ID:         activityID,
				Status:     &model.ActivityStatus{ID: getInt32FromRecord(record, "statusId")},
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
				// TODO: Populate User and Media
			}
			activities = append(activities, activity)
		}

		return activities, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.UserActivity), nil
}

// GetMediaActivities retrieves all activities for a media item
func (r *Neo4jRepository) GetMediaActivities(ctx context.Context, mediaID uuid.UUID) ([]*model.UserActivity, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:UserActivity)-[:ACTIVITY_FOR]->(m:Media {id: $mediaID})
			OPTIONAL MATCH (u:User)-[:HAS_ACTIVITY]->(a)
			RETURN a.id as id, a.statusId as statusId, a.rating as rating,
			       a.review as review, a.startedAt as startedAt, a.finishedAt as finishedAt,
			       u.id as userId
			ORDER BY a.createdAt DESC
		`

		params := map[string]any{"mediaID": mediaID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var activities []*model.UserActivity
		for result.Next(ctx) {
			record := result.Record()
			activityID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			activity := &model.UserActivity{
				ID:         activityID,
				Status:     &model.ActivityStatus{ID: getInt32FromRecord(record, "statusId")},
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
				// TODO: Populate User and Media
			}
			activities = append(activities, activity)
		}

		return activities, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.UserActivity), nil
}

// UpdateActivity updates an existing activity
func (r *Neo4jRepository) UpdateActivity(ctx context.Context, id uuid.UUID, statusID *int32, rating *float64, review *string, finishedAt *string) (*model.UserActivity, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:UserActivity {id: $id})
			SET a.updatedAt = datetime()
		`

		params := map[string]any{"id": id.String()}

		// Add optional fields to SET clause
		if statusID != nil {
			query += ", a.statusId = $statusId"
			params["statusId"] = *statusID
		}

		if rating != nil {
			query += ", a.rating = $rating"
			params["rating"] = *rating
		}

		if review != nil {
			query += ", a.review = $review"
			params["review"] = *review
		}

		if finishedAt != nil {
			query += ", a.finishedAt = $finishedAt"
			params["finishedAt"] = *finishedAt
		}

		query += `
			RETURN a.id as id, a.statusId as statusId, a.rating as rating,
			       a.review as review, a.startedAt as startedAt, a.finishedAt as finishedAt
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			activity := &model.UserActivity{
				ID:         id,
				Status:     &model.ActivityStatus{ID: getInt32FromRecord(record, "statusId")},
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
			}
			return activity, nil
		}

		return nil, fmt.Errorf("activity not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.UserActivity), nil
}

// DeleteActivity deletes an activity
func (r *Neo4jRepository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:UserActivity {id: $id})
			DETACH DELETE a
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

// Helper functions
func getFloat64Pointer(value interface{}) *float64 {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case float64:
		return &v
	case int32:
		f64 := float64(v)
		return &f64
	case int64:
		f64 := float64(v)
		return &f64
	}

	return nil
}

func getInt32FromRecord(record *neo4j.Record, key string) int32 {
	value := record.AsMap()[key]
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	}

	return 0
}
