package db

import (
	"context"
	"fmt"
	"nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateUser creates a new user in the database
func (r *Neo4jRepository) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	userID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (u:User {
				id: $id,
				name: $name,
				email: $email,
				authProvider: $authProvider,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN u.id as id, u.name as name, u.email as email, u.authProvider as authProvider
		`

		params := map[string]any{
			"id":           userID.String(),
			"name":         input.Name,
			"email":        input.Email,
			"authProvider": input.AuthProvider,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			user := &model.User{
				ID:              userID,
				Name:            record.AsMap()["name"].(string),
				Email:           record.AsMap()["email"].(string),
				AuthProvider:    getStringPointer(record.AsMap()["authProvider"]),
				Activities:      []*model.UserActivity{},
				Ratings:         []*model.Rating{},
				Favorites:       []model.Media{},
				Recommendations: []*model.Recommendation{},
			}
			return user, nil
		}

		return nil, fmt.Errorf("failed to create user")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

// GetUserByID retrieves a user by their ID
func (r *Neo4jRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $id})
			OPTIONAL MATCH (u)-[:HAS_ACTIVITY]->(a:UserActivity)
			OPTIONAL MATCH (u)-[:RATED]->(r:Rating)
			OPTIONAL MATCH (u)-[:FAVORITES]->(f:Media)
			OPTIONAL MATCH (u)-[:RECEIVED_RECOMMENDATION]->(rec:Recommendation)
			RETURN u.id as id, u.name as name, u.email as email, u.authProvider as authProvider,
			       collect(DISTINCT a) as activities,
			       collect(DISTINCT r) as ratings,
			       collect(DISTINCT f) as favorites,
			       collect(DISTINCT rec) as recommendations
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			user := &model.User{
				ID:              id,
				Name:            record.AsMap()["name"].(string),
				Email:           record.AsMap()["email"].(string),
				AuthProvider:    getStringPointer(record.AsMap()["authProvider"]),
				Activities:      []*model.UserActivity{},   // TODO: Parse activities
				Ratings:         []*model.Rating{},         // TODO: Parse ratings
				Favorites:       []model.Media{},           // TODO: Parse favorites
				Recommendations: []*model.Recommendation{}, // TODO: Parse recommendations
			}
			return user, nil
		}

		return nil, fmt.Errorf("user not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

// GetUserByEmail retrieves a user by their email
func (r *Neo4jRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {email: $email})
			RETURN u.id as id, u.name as name, u.email as email, u.authProvider as authProvider
		`

		params := map[string]any{"email": email}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			userID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			user := &model.User{
				ID:              userID,
				Name:            record.AsMap()["name"].(string),
				Email:           record.AsMap()["email"].(string),
				AuthProvider:    getStringPointer(record.AsMap()["authProvider"]),
				Activities:      []*model.UserActivity{},
				Ratings:         []*model.Rating{},
				Favorites:       []model.Media{},
				Recommendations: []*model.Recommendation{},
			}
			return user, nil
		}

		return nil, fmt.Errorf("user not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

// GetAllUsers retrieves all users
func (r *Neo4jRepository) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User)
			RETURN u.id as id, u.name as name, u.email as email, u.authProvider as authProvider
			ORDER BY u.name
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var users []*model.User
		for result.Next(ctx) {
			record := result.Record()
			userID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			user := &model.User{
				ID:              userID,
				Name:            record.AsMap()["name"].(string),
				Email:           record.AsMap()["email"].(string),
				AuthProvider:    getStringPointer(record.AsMap()["authProvider"]),
				Activities:      []*model.UserActivity{},
				Ratings:         []*model.Rating{},
				Favorites:       []model.Media{},
				Recommendations: []*model.Recommendation{},
			}
			users = append(users, user)
		}

		return users, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.User), nil
}

// UpdateUser updates an existing user
func (r *Neo4jRepository) UpdateUser(ctx context.Context, id uuid.UUID, input model.UpdateUserInput) (*model.User, error) {
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $id})
			SET u.updatedAt = datetime()
		`

		params := map[string]any{"id": id.String()}

		// Add optional fields to SET clause
		if input.Name != nil {
			query += ", u.name = $name"
			params["name"] = *input.Name
		}

		if input.Email != nil {
			query += ", u.email = $email"
			params["email"] = *input.Email
		}

		query += `
			RETURN u.id as id, u.name as name, u.email as email, u.authProvider as authProvider
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			user := &model.User{
				ID:              id,
				Name:            record.AsMap()["name"].(string),
				Email:           record.AsMap()["email"].(string),
				AuthProvider:    getStringPointer(record.AsMap()["authProvider"]),
				Activities:      []*model.UserActivity{},
				Ratings:         []*model.Rating{},
				Favorites:       []model.Media{},
				Recommendations: []*model.Recommendation{},
			}
			return user, nil
		}

		return nil, fmt.Errorf("user not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.User), nil
}

// DeleteUser deletes a user and all related data
func (r *Neo4jRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $id})
			DETACH DELETE u
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

// Helper function to safely get string pointer from interface{}
func getStringPointer(value interface{}) *string {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok && str != "" {
		return &str
	}
	return nil
}
