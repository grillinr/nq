package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Database wraps the Neo4j driver and provides database operations
type Database struct {
	driver neo4j.DriverWithContext
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	ctx := context.Background()

	// Load Neo4j environment variables
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

// validateAuraURI checks if the URI is properly formatted for Neo4j Aura
func validateAuraURI(uri string) error {
	if strings.HasPrefix(uri, "neo4j+s://") {
		if !strings.Contains(uri, "databases.neo4j.io") {
			return fmt.Errorf("invalid Aura URI format: should contain 'databases.neo4j.io'")
		}
		return nil
	} else if strings.HasPrefix(uri, "neo4j://") {
		if !strings.Contains(uri, "databases.neo4j.io") {
			return fmt.Errorf("invalid Aura URI format: should contain 'databases.neo4j.io'")
		}
		return nil
	} else if strings.HasPrefix(uri, "bolt://") {
		log.Println("Warning: Using bolt:// protocol. For Aura, consider using neo4j+s:// for better security")
		return nil
	}

	return fmt.Errorf("unsupported URI protocol. Use neo4j+s:// for Aura or bolt:// for local")
}

// Legacy function for backward compatibility
func DBConnect() neo4j.DriverWithContext {
	db, err := NewDatabase()
	if err != nil {
		panic(err)
	}
	return db.driver
}
