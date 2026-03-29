package db

import (
	"context"
	"fmt"
	"time"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// SendFriendRequest creates a FRIEND_REQUEST relationship from fromID to toID.
// Returns an error if a pending request or existing friendship already exists.
func (r *Neo4jRepository) SendFriendRequest(ctx context.Context, fromID, toID uuid.UUID) (*model.FriendRequest, error) {
	if fromID == toID {
		return nil, fmt.Errorf("cannot send a friend request to yourself")
	}

	requestID := uuid.New()
	createdAt := time.Now().UTC().Format(time.RFC3339)

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Atomically verify no existing friendship and MERGE the request to prevent race conditions.
		// If both users exist and are not friends, MERGE finds or creates the FRIEND_REQUEST.
		// ON CREATE SET only fires when the relationship is new; if it already exists wasCreated=false.
		query := `
			MATCH (from:User {id: $fromID}), (to:User {id: $toID})
			OPTIONAL MATCH (from)-[friendship:FRIEND_OF]->(to)
			WITH from, to, friendship IS NOT NULL as areFriends
			WHERE NOT areFriends
			MERGE (from)-[req:FRIEND_REQUEST]->(to)
			ON CREATE SET req.id = $requestID, req.createdAt = $createdAt
			RETURN req.id as requestId, req.createdAt as createdAt, req.id = $requestID as wasCreated
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"fromID":    fromID.String(),
			"toID":      toID.String(),
			"requestID": requestID.String(),
			"createdAt": createdAt,
		})
		if err != nil {
			return nil, err
		}
		if !result.Next(ctx) {
			return nil, fmt.Errorf("cannot send friend request: users not found or already friends")
		}
		row := result.Record().AsMap()
		if wasCreated, _ := row["wasCreated"].(bool); !wasCreated {
			return nil, fmt.Errorf("friend request already pending")
		}
		return requestID, nil
	})
	if err != nil {
		return nil, err
	}
	_ = result

	fromUser, err := r.GetUserByID(ctx, fromID)
	if err != nil {
		return nil, err
	}
	toUser, err := r.GetUserByID(ctx, toID)
	if err != nil {
		return nil, err
	}

	return &model.FriendRequest{
		ID:        requestID,
		From:      fromUser,
		To:        toUser,
		CreatedAt: createdAt,
	}, nil
}

// AcceptFriendRequest accepts a pending FRIEND_REQUEST by its ID.
// Creates bidirectional FRIEND_OF relationships and deletes the request.
// Returns the new friend (the user who sent the request).
func (r *Neo4jRepository) AcceptFriendRequest(ctx context.Context, requestID, acceptingUserID uuid.UUID) (*model.User, error) {
	var senderID uuid.UUID

	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Find the request where acceptingUser is the recipient
		findQuery := `
			MATCH (from:User)-[req:FRIEND_REQUEST {id: $requestID}]->(to:User {id: $acceptingUserID})
			RETURN from.id as fromId
		`
		findResult, err := tx.Run(ctx, findQuery, map[string]any{
			"requestID":       requestID.String(),
			"acceptingUserID": acceptingUserID.String(),
		})
		if err != nil {
			return nil, err
		}
		if !findResult.Next(ctx) {
			return nil, fmt.Errorf("friend request not found or you are not the recipient")
		}
		fromIDStr, _ := findResult.Record().AsMap()["fromId"].(string)
		parsed, err := uuid.Parse(fromIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid sender ID in request: %w", err)
		}
		senderID = parsed

		since := time.Now().UTC().Format(time.RFC3339)

		// Create bidirectional FRIEND_OF and delete the request
		acceptQuery := `
			MATCH (from:User {id: $fromID})-[req:FRIEND_REQUEST {id: $requestID}]->(to:User {id: $toID})
			DELETE req
			MERGE (from)-[:FRIEND_OF {since: $since}]->(to)
			MERGE (to)-[:FRIEND_OF {since: $since}]->(from)
		`
		_, err = tx.Run(ctx, acceptQuery, map[string]any{
			"requestID": requestID.String(),
			"fromID":    senderID.String(),
			"toID":      acceptingUserID.String(),
			"since":     since,
		})
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(ctx, senderID)
}

// DeclineFriendRequest removes a pending FRIEND_REQUEST relationship.
func (r *Neo4jRepository) DeclineFriendRequest(ctx context.Context, requestID, decliningUserID uuid.UUID) (bool, error) {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (from:User)-[req:FRIEND_REQUEST {id: $requestID}]->(to:User {id: $decliningUserID})
			DELETE req
			RETURN count(req) as deleted
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"requestID":       requestID.String(),
			"decliningUserID": decliningUserID.String(),
		})
		if err != nil {
			return nil, err
		}
		if !result.Next(ctx) {
			return nil, fmt.Errorf("friend request not found")
		}
		return nil, nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoveFriend removes the bidirectional FRIEND_OF relationship between two users.
func (r *Neo4jRepository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {id: $userID})-[r1:FRIEND_OF]->(b:User {id: $friendID})
			OPTIONAL MATCH (b)-[r2:FRIEND_OF]->(a)
			DELETE r1, r2
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"userID":   userID.String(),
			"friendID": friendID.String(),
		})
		return nil, err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetFriends returns all friends of a user (users connected via FRIEND_OF).
func (r *Neo4jRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userID})-[:FRIEND_OF]->(friend:User)
			RETURN friend.id as id, friend.name as name, friend.email as email,
			       friend.authProvider as authProvider, friend.authSubject as authSubject,
			       friend.avatarUrl as avatarUrl
			ORDER BY friend.name
		`
		records, err := tx.Run(ctx, query, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}
		var users []*model.User
		for records.Next(ctx) {
			u, err := userFromRecord(records.Record())
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
		if err := records.Err(); err != nil {
			return nil, err
		}
		return users, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []*model.User{}, nil
	}
	return result.([]*model.User), nil
}

// GetPendingFriendRequests returns all FRIEND_REQUEST relationships directed at userID (received requests).
func (r *Neo4jRepository) GetPendingFriendRequests(ctx context.Context, userID uuid.UUID) ([]*model.FriendRequest, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (from:User)-[req:FRIEND_REQUEST]->(to:User {id: $userID})
			RETURN req.id as requestId, req.createdAt as createdAt,
			       from.id as fromId, from.name as fromName, from.email as fromEmail,
			       from.authProvider as fromAuthProvider, from.authSubject as fromAuthSubject,
			       from.avatarUrl as fromAvatarUrl,
			       to.id as toId, to.name as toName, to.email as toEmail,
			       to.authProvider as toAuthProvider, to.authSubject as toAuthSubject,
			       to.avatarUrl as toAvatarUrl
			ORDER BY req.createdAt DESC
		`
		records, err := tx.Run(ctx, query, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}
		var requests []*model.FriendRequest
		for records.Next(ctx) {
			req, err := friendRequestFromRecord(records.Record())
			if err != nil {
				return nil, err
			}
			requests = append(requests, req)
		}
		if err := records.Err(); err != nil {
			return nil, err
		}
		return requests, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []*model.FriendRequest{}, nil
	}
	return result.([]*model.FriendRequest), nil
}

// GetSentFriendRequests returns all FRIEND_REQUEST relationships sent by userID.
func (r *Neo4jRepository) GetSentFriendRequests(ctx context.Context, userID uuid.UUID) ([]*model.FriendRequest, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (from:User {id: $userID})-[req:FRIEND_REQUEST]->(to:User)
			RETURN req.id as requestId, req.createdAt as createdAt,
			       from.id as fromId, from.name as fromName, from.email as fromEmail,
			       from.authProvider as fromAuthProvider, from.authSubject as fromAuthSubject,
			       from.avatarUrl as fromAvatarUrl,
			       to.id as toId, to.name as toName, to.email as toEmail,
			       to.authProvider as toAuthProvider, to.authSubject as toAuthSubject,
			       to.avatarUrl as toAvatarUrl
			ORDER BY req.createdAt DESC
		`
		records, err := tx.Run(ctx, query, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}
		var requests []*model.FriendRequest
		for records.Next(ctx) {
			req, err := friendRequestFromRecord(records.Record())
			if err != nil {
				return nil, err
			}
			requests = append(requests, req)
		}
		if err := records.Err(); err != nil {
			return nil, err
		}
		return requests, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []*model.FriendRequest{}, nil
	}
	return result.([]*model.FriendRequest), nil
}

// SearchUsers searches for users by name (case-insensitive contains), excluding the given user.
func (r *Neo4jRepository) SearchUsers(ctx context.Context, query string, excludeUserID uuid.UUID) ([]*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `
			MATCH (u:User)
			WHERE u.id <> $excludeID
			  AND toLower(u.name) CONTAINS toLower($query)
			RETURN u.id as id, u.name as name, u.email as email,
			       u.authProvider as authProvider, u.authSubject as authSubject,
			       u.avatarUrl as avatarUrl
			ORDER BY u.name
			LIMIT 20
		`
		records, err := tx.Run(ctx, cypher, map[string]any{
			"query":     query,
			"excludeID": excludeUserID.String(),
		})
		if err != nil {
			return nil, err
		}
		var users []*model.User
		for records.Next(ctx) {
			u, err := userFromRecord(records.Record())
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
		if err := records.Err(); err != nil {
			return nil, err
		}
		return users, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []*model.User{}, nil
	}
	return result.([]*model.User), nil
}

// GetFriendsActivity returns recent activities from all friends of the given user.
func (r *Neo4jRepository) GetFriendsActivity(ctx context.Context, userID uuid.UUID, limit int) ([]*model.UserActivity, error) {
	if limit <= 0 {
		limit = 50
	}

	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (me:User {id: $userID})-[:FRIEND_OF]->(friend:User)
			MATCH (friend)-[ha:HAS_ACTIVITY]->(m:Media)
			RETURN ha.id as id, ha.statusId as statusId, ha.rating as rating,
			       ha.review as review, ha.startedAt as startedAt, ha.finishedAt as finishedAt,
			       friend.id as userId, m.id as mediaId
			ORDER BY ha.updatedAt DESC, ha.createdAt DESC
			LIMIT $limit
		`
		records, err := tx.Run(ctx, query, map[string]any{
			"userID": userID.String(),
			"limit":  int64(limit),
		})
		if err != nil {
			return nil, err
		}

		var activities []*model.UserActivity
		for records.Next(ctx) {
			record := records.Record()
			m := record.AsMap()

			idStr, _ := m["id"].(string)
			activityID, err := uuid.Parse(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid activity id: %w", err)
			}

			activity := &model.UserActivity{
				ID:         activityID,
				Status:     &model.ActivityStatus{ID: getInt32FromRecord(record, "statusId")},
				Rating:     getFloat64Pointer(m["rating"]),
				Review:     getStringPointer(m["review"]),
				StartedAt:  getStringPointer(m["startedAt"]),
				FinishedAt: getStringPointer(m["finishedAt"]),
			}

			// Populate full status name
			if statusID := getInt32FromRecord(record, "statusId"); statusID > 0 {
				if s, err := r.GetActivityStatusByID(ctx, statusID); err == nil && s != nil {
					activity.Status = s
				}
			}

			// Populate friend user
			if userIDStr, ok := m["userId"].(string); ok && userIDStr != "" {
				if uid, err := uuid.Parse(userIDStr); err == nil {
					if u, err := r.GetUserByID(ctx, uid); err == nil {
						activity.User = u
					}
				}
			}

			// Populate media
			if mediaIDStr, ok := m["mediaId"].(string); ok && mediaIDStr != "" {
				if mid, err := uuid.Parse(mediaIDStr); err == nil {
					if media, err := r.GetMediaByID(ctx, mid); err == nil {
						activity.Media = media
					}
				}
			}

			activities = append(activities, activity)
		}
		if err := records.Err(); err != nil {
			return nil, err
		}
		return activities, nil
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []*model.UserActivity{}, nil
	}
	return result.([]*model.UserActivity), nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func userFromRecord(record *neo4j.Record) (*model.User, error) {
	m := record.AsMap()
	idStr, _ := m["id"].(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %q: %w", idStr, err)
	}
	return &model.User{
		ID:                    id,
		Name:                  m["name"].(string),
		Email:                 m["email"].(string),
		AuthProvider:          getStringPointer(m["authProvider"]),
		AuthSubject:           getStringPointer(m["authSubject"]),
		AvatarURL:             getStringPointer(m["avatarUrl"]),
		Activities:            []*model.UserActivity{},
		Ratings:               []*model.Rating{},
		Favorites:             []model.Media{},
		Recommendations:       []*model.Recommendation{},
		Friends:               []*model.User{},
		PendingFriendRequests: []*model.FriendRequest{},
		SentFriendRequests:    []*model.FriendRequest{},
	}, nil
}

func friendRequestFromRecord(record *neo4j.Record) (*model.FriendRequest, error) {
	m := record.AsMap()

	requestIDStr, _ := m["requestId"].(string)
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid request id: %w", err)
	}

	createdAt, _ := m["createdAt"].(string)

	fromIDStr, _ := m["fromId"].(string)
	fromID, err := uuid.Parse(fromIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid from user id: %w", err)
	}
	fromUser := &model.User{
		ID:                    fromID,
		Name:                  stringFromMap(m, "fromName"),
		Email:                 stringFromMap(m, "fromEmail"),
		AuthProvider:          getStringPointer(m["fromAuthProvider"]),
		AuthSubject:           getStringPointer(m["fromAuthSubject"]),
		AvatarURL:             getStringPointer(m["fromAvatarUrl"]),
		Activities:            []*model.UserActivity{},
		Ratings:               []*model.Rating{},
		Favorites:             []model.Media{},
		Recommendations:       []*model.Recommendation{},
		Friends:               []*model.User{},
		PendingFriendRequests: []*model.FriendRequest{},
		SentFriendRequests:    []*model.FriendRequest{},
	}

	toIDStr, _ := m["toId"].(string)
	toID, err := uuid.Parse(toIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid to user id: %w", err)
	}
	toUser := &model.User{
		ID:                    toID,
		Name:                  stringFromMap(m, "toName"),
		Email:                 stringFromMap(m, "toEmail"),
		AuthProvider:          getStringPointer(m["toAuthProvider"]),
		AuthSubject:           getStringPointer(m["toAuthSubject"]),
		AvatarURL:             getStringPointer(m["toAvatarUrl"]),
		Activities:            []*model.UserActivity{},
		Ratings:               []*model.Rating{},
		Favorites:             []model.Media{},
		Recommendations:       []*model.Recommendation{},
		Friends:               []*model.User{},
		PendingFriendRequests: []*model.FriendRequest{},
		SentFriendRequests:    []*model.FriendRequest{},
	}

	return &model.FriendRequest{
		ID:        requestID,
		From:      fromUser,
		To:        toUser,
		CreatedAt: createdAt,
	}, nil
}

func stringFromMap(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
