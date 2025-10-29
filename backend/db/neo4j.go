// Package db provides a wrapper around the Neo4j driver for database operations
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Database wraps the Neo4j driver and provides database operations
type Database struct {
	driver neo4j.DriverWithContext
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	ctx := context.Background()

	// Load Neo4j environment variables (also try loading a .env from repo root)
	_ = loadEnvUpwards(".env", 4)
	dbURI := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USERNAME")
	dbPassword := os.Getenv("NEO4J_PASSWORD")

	// Create driver connection (Aura handles optimization automatically)
	driver, err := neo4j.NewDriverWithContext(
		dbURI,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connectivity
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		driver.Close(ctx)
		return nil, fmt.Errorf("failed to verify Neo4j connectivity: %w", err)
	}

	log.Println("Successfully connected to Neo4j Aura")

	return &Database{driver: driver}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	ctx := context.Background()
	return db.driver.Close(ctx)
}

// GetDriver returns the underlying Neo4j driver
func (db *Database) GetDriver() neo4j.DriverWithContext {
	return db.driver
}

// ExecuteRead executes a read transaction
func (db *Database) ExecuteRead(ctx context.Context, work func(neo4j.ManagedTransaction) (any, error)) (any, error) {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, work)
	return result, err
}

// ExecuteWrite executes a write transaction
func (db *Database) ExecuteWrite(ctx context.Context, work func(neo4j.ManagedTransaction) (any, error)) (any, error) {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, work)
	return result, err
}
