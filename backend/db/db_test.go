package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TestDatabase wraps the regular Database for testing
type TestDatabase struct {
	*Database
	isTestDB bool
}

// NewTestDatabase creates a test database connection
func NewTestDatabase(t *testing.T) *TestDatabase {
	// Skip tests if no test database is available
	if os.Getenv("NEO4J_TEST_URI") == "" {
		t.Skip("Skipping database tests: NEO4J_TEST_URI not set")
	}

	ctx := context.Background()

	// Use separate test database environment variables
	dbURI := os.Getenv("NEO4J_TEST_URI")
	dbUser := os.Getenv("NEO4J_TEST_USERNAME")
	dbPassword := os.Getenv("NEO4J_TEST_PASSWORD")

	if dbURI == "" {
		dbURI = os.Getenv("NEO4J_URI") // Fallback to main DB for local testing
	}
	if dbUser == "" {
		dbUser = os.Getenv("NEO4J_USERNAME")
	}
	if dbPassword == "" {
		dbPassword = os.Getenv("NEO4J_PASSWORD")
	}

	// Create driver connection
	driver, err := neo4j.NewDriverWithContext(
		dbURI,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		t.Fatalf("Failed to create test Neo4j driver: %v", err)
	}

	// Verify connectivity
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		driver.Close(ctx)
		t.Fatalf("Failed to verify test Neo4j connectivity: %v", err)
	}

	db := &Database{driver: driver}
	testDB := &TestDatabase{Database: db, isTestDB: true}

	// Clean up any existing test data
	testDB.CleanupTestData(t)

	return testDB
}

// CleanupTestData removes all test data from the database
func (tdb *TestDatabase) CleanupTestData(t *testing.T) {
	ctx := context.Background()

	// Delete all nodes and relationships with test prefix or in test session
	_, err := tdb.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Delete test data - be very careful with this query
		queries := []string{
			"MATCH (n) WHERE n.id STARTS WITH 'test-' DETACH DELETE n",
			"MATCH (n) WHERE n.name STARTS WITH 'test-' DETACH DELETE n",
			"MATCH (n) WHERE n.email STARTS WITH 'test-' DETACH DELETE n",
			"MATCH (n) WHERE n.title STARTS WITH 'test-' DETACH DELETE n",
		}

		for _, query := range queries {
			_, err := tx.Run(ctx, query, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to cleanup test data: %w", err)
			}
		}
		return nil, nil
	})

	if err != nil {
		t.Logf("Warning: Failed to cleanup test data: %v", err)
	}
}

// Close closes the test database connection and cleans up
func (tdb *TestDatabase) Close(t *testing.T) {
	if tdb.isTestDB {
		tdb.CleanupTestData(t)
	}
	tdb.Database.Close()
}

// GenerateTestID generates a UUID with test prefix for easy cleanup
func GenerateTestID() uuid.UUID {
	id := uuid.New()
	// Use a deterministic test prefix that we can clean up later
	testStr := fmt.Sprintf("test-%s", id.String()[5:]) // Keep most of the UUID but add test prefix
	return uuid.MustParse(fmt.Sprintf("00000000-%s", testStr[5:]))
}

// GenerateTestEmail generates a test email address
func GenerateTestEmail() string {
	return fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8])
}

// GenerateTestName generates a test name
func GenerateTestName() string {
	return fmt.Sprintf("test-user-%s", uuid.New().String()[:8])
}

// GenerateTestTitle generates a test title
func GenerateTestTitle(mediaType string) string {
	return fmt.Sprintf("test-%s-%s", mediaType, uuid.New().String()[:8])
}

// setupTestRepository creates a test repository with a test database
func setupTestRepository(t *testing.T) (*Neo4jRepository, *TestDatabase) {
	testDB := NewTestDatabase(t)
	repo := NewNeo4jRepository(testDB.Database)
	return repo, testDB
}

// TestMain handles setup and teardown for all tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Exit with the same code as the test run
	os.Exit(code)
}
