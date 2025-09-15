package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateConstraints creates all necessary Neo4j constraints and indexes
func (db *Database) CreateConstraints(ctx context.Context) error {
	constraints := []string{
		// User constraints
		"CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE",
		"CREATE CONSTRAINT user_email_unique IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE",

		// Media constraints
		"CREATE CONSTRAINT media_id_unique IF NOT EXISTS FOR (m:Media) REQUIRE m.id IS UNIQUE",
		"CREATE CONSTRAINT movie_id_unique IF NOT EXISTS FOR (m:Movie) REQUIRE m.id IS UNIQUE",
		"CREATE CONSTRAINT tvshow_id_unique IF NOT EXISTS FOR (t:TVShow) REQUIRE t.id IS UNIQUE",
		"CREATE CONSTRAINT book_id_unique IF NOT EXISTS FOR (b:Book) REQUIRE b.id IS UNIQUE",
		"CREATE CONSTRAINT game_id_unique IF NOT EXISTS FOR (g:Game) REQUIRE g.id IS UNIQUE",
		"CREATE CONSTRAINT musicalbum_id_unique IF NOT EXISTS FOR (ma:MusicAlbum) REQUIRE ma.id IS UNIQUE",

		// Creator constraints
		"CREATE CONSTRAINT creator_id_unique IF NOT EXISTS FOR (c:Creator) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT creatorrole_id_unique IF NOT EXISTS FOR (cr:CreatorRole) REQUIRE cr.id IS UNIQUE",

		// Platform constraints
		"CREATE CONSTRAINT platform_id_unique IF NOT EXISTS FOR (p:Platform) REQUIRE p.id IS UNIQUE",

		// Activity constraints
		"CREATE CONSTRAINT activity_id_unique IF NOT EXISTS FOR (a:UserActivity) REQUIRE a.id IS UNIQUE",
		"CREATE CONSTRAINT activitystatus_id_unique IF NOT EXISTS FOR (as:ActivityStatus) REQUIRE as.id IS UNIQUE",

		// Tag constraints
		"CREATE CONSTRAINT tag_id_unique IF NOT EXISTS FOR (t:Tag) REQUIRE t.id IS UNIQUE",

		// Rating constraints - composite unique constraint for user+media
		"CREATE CONSTRAINT rating_user_media_unique IF NOT EXISTS FOR (r:Rating) REQUIRE (r.userId, r.mediaId) IS UNIQUE",

		// Recommendation constraints
		"CREATE CONSTRAINT recommendation_id_unique IF NOT EXISTS FOR (r:Recommendation) REQUIRE r.id IS UNIQUE",
	}

	for _, constraint := range constraints {
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, constraint, nil)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})

		if err != nil {
			return fmt.Errorf("failed to create constraint '%s': %w", constraint, err)
		}
	}

	return nil
}

// CreateIndexes creates all necessary Neo4j indexes for better performance
func (db *Database) CreateIndexes(ctx context.Context) error {
	indexes := []string{
		// Media indexes
		"CREATE INDEX media_title_index IF NOT EXISTS FOR (m:Media) ON (m.title)",
		"CREATE INDEX media_release_date_index IF NOT EXISTS FOR (m:Media) ON (m.releaseDate)",
		"CREATE INDEX movie_title_index IF NOT EXISTS FOR (m:Movie) ON (m.title)",
		"CREATE INDEX tvshow_title_index IF NOT EXISTS FOR (t:TVShow) ON (t.title)",
		"CREATE INDEX book_title_index IF NOT EXISTS FOR (b:Book) ON (b.title)",
		"CREATE INDEX game_title_index IF NOT EXISTS FOR (g:Game) ON (g.title)",
		"CREATE INDEX musicalbum_title_index IF NOT EXISTS FOR (ma:MusicAlbum) ON (ma.title)",

		// User indexes
		"CREATE INDEX user_name_index IF NOT EXISTS FOR (u:User) ON (u.name)",
		"CREATE INDEX user_email_index IF NOT EXISTS FOR (u:User) ON (u.email)",

		// Creator indexes
		"CREATE INDEX creator_name_index IF NOT EXISTS FOR (c:Creator) ON (c.name)",

		// Platform indexes
		"CREATE INDEX platform_name_index IF NOT EXISTS FOR (p:Platform) ON (p.name)",

		// Tag indexes
		"CREATE INDEX tag_name_index IF NOT EXISTS FOR (t:Tag) ON (t.name)",
		"CREATE INDEX tag_type_index IF NOT EXISTS FOR (t:Tag) ON (t.type)",

		// Activity indexes
		"CREATE INDEX activity_user_index IF NOT EXISTS FOR (a:UserActivity) ON (a.userId)",
		"CREATE INDEX activity_media_index IF NOT EXISTS FOR (a:UserActivity) ON (a.mediaId)",
		"CREATE INDEX activity_status_index IF NOT EXISTS FOR (a:UserActivity) ON (a.statusId)",

		// Rating indexes
		"CREATE INDEX rating_user_index IF NOT EXISTS FOR (r:Rating) ON (r.userId)",
		"CREATE INDEX rating_media_index IF NOT EXISTS FOR (r:Rating) ON (r.mediaId)",
		"CREATE INDEX rating_score_index IF NOT EXISTS FOR (r:Rating) ON (r.score)",

		// Recommendation indexes
		"CREATE INDEX recommendation_user_index IF NOT EXISTS FOR (r:Recommendation) ON (r.userId)",
		"CREATE INDEX recommendation_media_index IF NOT EXISTS FOR (r:Recommendation) ON (r.mediaId)",
	}

	for _, index := range indexes {
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, index, nil)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})

		if err != nil {
			return fmt.Errorf("failed to create index '%s': %w", index, err)
		}
	}

	return nil
}

// InitializeDatabase creates all constraints and indexes
func (db *Database) InitializeDatabase(ctx context.Context) error {
	if err := db.CreateConstraints(ctx); err != nil {
		return fmt.Errorf("failed to create constraints: %w", err)
	}

	if err := db.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
