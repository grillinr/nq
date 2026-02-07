package db

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// BackfillPersonIDs assigns missing ids on Person nodes.
func (db *Database) BackfillPersonIDs(ctx context.Context) error {
	_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Person)
			WHERE p.id IS NULL
			SET p.id = randomUUID()
			RETURN count(p) as updated
		`
		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}
		if _, err := result.Consume(ctx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("backfill person ids failed: %w", err)
	}
	return nil
}
