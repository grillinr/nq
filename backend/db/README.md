# Database

This package contains the database layer for the NQ application, implementing repositories for users, media, activities, ratings, and recommendations using Neo4j as the graph database.

## Structure

- `neo4j.go` - Main database connection and configuration
- `repositories.go` - Repository interface definitions
- `user_repository.go` - User-related database operations
- `media_repository.go` - Media (movies, books, games, etc.) operations
- `activity_repository.go` - User activity tracking
- `rating_repository.go` - User ratings for media
- `recommendation_repository.go` - Recommendation system operations
- `constraints.go` - Database constraints and indexes

## Usage

The repository pattern is used to abstract database operations. Each repository implements a specific interface and provides methods for CRUD operations on their respective entities.

### Example

```go
// Initialize repository
repo := NewNeo4jRepository(driver)

// Create a user
user, err := repo.CreateUser(ctx, "john@example.com", "John Doe", "hashed_password")

// Get user by ID
user, err := repo.GetUserByID(ctx, userID)
```

## Testing

This directory contains comprehensive tests for all repository operations. See the testing section below for details.

### Setup

#### Prerequisites

1. **Neo4j Database**: You need a running Neo4j database instance for testing
2. **Environment Variables**: Set up the required environment variables

#### Environment Variables

Create a `.env` file in the backend directory or set these environment variables:

```bash
# Required for tests to run (otherwise they will be skipped)
NEO4J_TEST_URI=neo4j://localhost:7687

# Neo4j credentials
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_password_here
```

#### Running Tests

```bash
# Run all database tests
go test ./db -v

# Run specific test files
go test ./db/user_repository_test.go -v
go test ./db/media_repository_test.go -v
go test ./db/activity_repository_test.go -v
go test ./db/rating_repository_test.go -v
go test ./db/recommendation_repository_test.go -v

# Run with coverage
go test ./db -cover

# Run specific test function
go test ./db -run TestUserRepository_CreateUser -v
```

### Test Structure

#### Test Files

- `db_test.go` - Test infrastructure and utilities
- `user_repository_test.go` - User CRUD operations (7 test functions)
- `media_repository_test.go` - Media operations for all media types (9 test functions)
- `activity_repository_test.go` - User activity tracking (6 test functions)
- `rating_repository_test.go` - User ratings for media (7 test functions)
- `recommendation_repository_test.go` - Recommendation system tests (4 test functions)

#### Test Infrastructure

The tests use a shared infrastructure defined in `db_test.go`:

- **TestDatabase**: Wrapper around Neo4j connection with cleanup utilities
- **Test Data Generators**: Helper functions for creating test users, media, etc.
- **Cleanup**: Automatic cleanup after each test to ensure isolation

#### Key Features

1. **Isolation**: Each test runs in isolation with its own test data
2. **Cleanup**: Automatic cleanup prevents test interference
3. **Realistic Data**: Uses realistic test data with proper relationships
4. **Comprehensive Coverage**: Tests all CRUD operations and edge cases
5. **Error Handling**: Tests both success and failure scenarios
6. **GraphQL Compatibility**: Tests work with generated GraphQL models

### Test Coverage

#### User Repository Tests ✅

- User creation and validation
- User retrieval by ID and email
- User updates and deletion
- Favorites management
- Error handling for invalid data

#### Media Repository Tests ✅

- Movie creation and retrieval
- Book creation with ISBN validation
- Game creation with platform support
- TV Show creation with episode handling
- Music Album creation with artist info
- Generic media operations

#### Activity Repository Tests ✅

- Activity creation and tracking
- Status updates (watching, completed, etc.)
- User and media activity queries
- Progress tracking
- GraphQL model compatibility

#### Rating Repository Tests ✅

- Rating creation and updates
- User and media rating queries
- Average rating calculations
- Rating deletion and validation

#### Recommendation Repository Tests ✅

- Recommendation creation with source and score
- Recommendation retrieval by ID
- User recommendations listing
- Recommendation deletion
