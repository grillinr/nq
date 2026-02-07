package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type activityStatusSeed struct {
	id   int32
	name string
}

var defaultActivityStatuses = []activityStatusSeed{
	{id: 1, name: "Planned"},
	{id: 2, name: "In Progress"},
	{id: 3, name: "Completed"},
}

// SeedActivityStatuses ensures default activity statuses exist.
func (db *Database) SeedActivityStatuses(ctx context.Context) error {
	for _, status := range defaultActivityStatuses {
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
				MERGE (s:ActivityStatus {id: $id})
				SET s.name = $name
				RETURN s.id as id
			`
			params := map[string]any{
				"id":   status.id,
				"name": status.name,
			}
			result, err := tx.Run(ctx, query, params)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})
		if err != nil {
			return fmt.Errorf("failed to seed activity status %d: %w", status.id, err)
		}
	}

	return nil
}
