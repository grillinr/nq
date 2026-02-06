package db

import (
	"context"
	"fmt"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateActivity creates a new user activity in the database
func (r *Neo4jRepository) CreateActivity(ctx context.Context, userID uuid.UUID, input model.CreateActivityInput) (*model.UserActivity, error) {
	activityID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Create a HAS_ACTIVITY relationship between User and Media and store activity properties on the relationship
		query := `
			MATCH (u:User {id: $userID}), (m:Media {id: $mediaID})
			CREATE (u)-[ha:HAS_ACTIVITY {
				id: $activityID,
				statusId: $statusID,
				rating: $rating,
				review: $review,
				startedAt: $startedAt,
				finishedAt: $finishedAt,
				createdAt: datetime(),
				updatedAt: datetime()
			}]->(m)
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
				   ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt
		`

		params := map[string]any{
			"activityID": activityID.String(),
			"userID":     userID.String(),
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
			// parse id from returned string
			idStr, _ := record.Get("id")
			activityID, _ := uuid.Parse(idStr.(string))

			activity := &model.UserActivity{
				ID:         activityID,
				Status:     &model.ActivityStatus{ID: input.StatusID},
				Rating:     getFloat64Pointer(record.AsMap()["rating"]),
				Review:     getStringPointer(record.AsMap()["review"]),
				StartedAt:  getStringPointer(record.AsMap()["startedAt"]),
				FinishedAt: getStringPointer(record.AsMap()["finishedAt"]),
			}

			// Try to populate full status, user and media if available. Failures are non-fatal.
			if s, err := r.GetActivityStatusByID(ctx, input.StatusID); err == nil && s != nil {
				activity.Status = s
			}
			if u, err := r.GetUserByID(ctx, userID); err == nil {
				activity.User = u
			}
			if m, err := r.GetMediaByID(ctx, input.MediaID); err == nil {
				activity.Media = m
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
		// Find the HAS_ACTIVITY relationship with matching id
		query := `
			MATCH (u:User)-[ha:HAS_ACTIVITY]->(m:Media)
			WHERE ha.id = $id
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
				   ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt,
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
			}
			// populate user and media if ids returned
			if userIDStr, ok := record.AsMap()["userId"].(string); ok && userIDStr != "" {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						activity.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok && mediaIDStr != "" {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						activity.Media = m
					}
				}
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
			MATCH (u:User {id: $userID})-[ha:HAS_ACTIVITY]->(m:Media)
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
				   ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt,
				   m.id as mediaId
			ORDER BY ha.createdAt DESC
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
			}
			// Populate media for this user's activity if mediaId returned
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok && mediaIDStr != "" {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						activity.Media = m
					}
				}
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
			MATCH (u:User)-[ha:HAS_ACTIVITY]->(m:Media {id: $mediaID})
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
				   ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt,
				   u.id as userId
			ORDER BY ha.createdAt DESC
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
			}
			// Populate user for this media activity if userId returned
			if userIDStr, ok := record.AsMap()["userId"].(string); ok && userIDStr != "" {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						activity.User = u
					}
				}
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
			MATCH (u:User)-[ha:HAS_ACTIVITY]->(m:Media)
			WHERE ha.id = $id
			SET ha.updatedAt = datetime()
		`

		params := map[string]any{"id": id.String()}

		// Add optional fields to SET clause
		if statusID != nil {
			query += ", ha.statusId = $statusId"
			params["statusId"] = *statusID
		}

		if rating != nil {
			query += ", ha.rating = $rating"
			params["rating"] = *rating
		}

		if review != nil {
			query += ", ha.review = $review"
			params["review"] = *review
		}

		if finishedAt != nil {
			query += ", ha.finishedAt = $finishedAt"
			params["finishedAt"] = *finishedAt
		}

		query += `
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
				   ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt,
				   u.id as userId, m.id as mediaId
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
			// try populate user and media
			if userIDStr, ok := record.AsMap()["userId"].(string); ok && userIDStr != "" {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						activity.User = u
					}
				}
			}
			if mediaIDStr, ok := record.AsMap()["mediaId"].(string); ok && mediaIDStr != "" {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if m, err := r.GetMediaByID(ctx, mid); err == nil {
						activity.Media = m
					}
				}
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
		// Delete the relationship that stores the activity
		query := `
			MATCH ()-[ha:HAS_ACTIVITY]->()
			WHERE ha.id = $id
			DELETE ha
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, err
		}

		// If no relationships were deleted, the activity did not exist
		if summary.Counters().RelationshipsDeleted() == 0 {
			return nil, fmt.Errorf("activity not found")
		}

		return summary, nil
	})

	return err
}

// Helper functions
func getFloat64Pointer(value any) *float64 {
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

// GetActivityStatusByID looks up an ActivityStatus node by its ID and returns a model.ActivityStatus.
// This keeps ActivityStatus resolution local and optional (errors returned if not found).
func (r *Neo4jRepository) GetActivityStatusByID(ctx context.Context, id int32) (*model.ActivityStatus, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (s:ActivityStatus {id: $id})
			RETURN s.id as id, s.name as name
		`
		params := map[string]any{"id": id}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			record := result.Record()
			status := &model.ActivityStatus{
				ID:   getInt32FromRecord(record, "id"),
				Name: record.AsMap()["name"].(string),
			}
			return status, nil
		}
		return nil, fmt.Errorf("activity status not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.ActivityStatus), nil
}
