# Database Layer

This directory contains the Neo4j database implementation for the NQ (Not Quite) media tracking application.

## Architecture

The database layer follows the Repository pattern with the following structure:

- **Database**: Main database connection wrapper
- **Repository Interface**: Defines all database operations
- **Neo4jRepository**: Concrete implementation using Neo4j
- **Specialized Repositories**: Separate files for different entity types

## Files

### Core Files
- `neo4j.go` - Database connection and session management
- `repositories.go` - Repository interfaces and main implementation
- `constraints.go` - Neo4j constraints and indexes setup

### Repository Implementations
- `user_repository.go` - User CRUD operations
- `media_repository.go` - Media (Movie, TV Show, Book, Game, Music Album) operations
- `activity_repository.go` - User activity tracking
- `rating_repository.go` - Rating system
- `recommendation_repository.go` - Recommendation engine

## Neo4j Schema

### Node Types
- **User**: Users of the application
- **Media**: Base interface for all media types
- **Movie**: Movies
- **TVShow**: Television shows
- **Book**: Books
- **Game**: Video games
- **MusicAlbum**: Music albums
- **Creator**: Media creators (directors, authors, etc.)
- **Platform**: Streaming platforms and stores
- **Tag**: Media tags and categories
- **UserActivity**: User interactions with media
- **Rating**: User ratings of media
- **Recommendation**: Media recommendations

### Relationships
- `(User)-[:HAS_ACTIVITY]->(UserActivity)`
- `(UserActivity)-[:ACTIVITY_FOR]->(Media)`
- `(User)-[:RATED]->(Rating)`
- `(Rating)-[:RATING_FOR]->(Media)`
- `(User)-[:FAVORITES]->(Media)`
- `(User)-[:RECEIVED_RECOMMENDATION]->(Recommendation)`
- `(Recommendation)-[:RECOMMENDS]->(Media)`
- `(Creator)-[:CREATED]->(Media)`
- `(Platform)-[:HOSTS]->(Media)`
- `(Media)-[:TAGGED_WITH]->(Tag)`

## Usage

### Initialization
```go
// Create database connection
db, err := db.NewDatabase()
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Initialize constraints and indexes
ctx := context.Background()
if err := db.InitializeDatabase(ctx); err != nil {
    log.Printf("Warning: %v", err)
}

// Create repository
repo := db.NewNeo4jRepository(db)
```

### Basic Operations
```go
// Create a user
user, err := repo.CreateUser(ctx, model.CreateUserInput{
    Name:  "John Doe",
    Email: "john@example.com",
})

// Get user by ID
user, err := repo.GetUserByID(ctx, userID)

// Create a movie
movie, err := repo.CreateMovie(ctx, model.CreateMovieInput{
    Title: "Inception",
    Description: "A mind-bending thriller",
})
```

## Environment Variables

Required environment variables for Neo4j Aura:
- `NEO4J_URI`: Neo4j Aura connection URI (format: neo4j+s://your-instance-id.databases.neo4j.io)
- `NEO4J_USERNAME`: Neo4j username (default: neo4j)
- `NEO4J_PASSWORD`: Neo4j Aura password (from your Aura dashboard)

### Neo4j Aura Setup

1. **Create Aura Instance**: Go to [Neo4j Aura](https://console.neo4j.io/) and create a new database
2. **Get Connection Details**: Copy the connection URI and password from your Aura dashboard
3. **Set Environment Variables**: Update your `.env` file with the Aura credentials

Example Aura URI format:
```
neo4j+s://abcd1234.databases.neo4j.io
```

**Note**: Aura uses `neo4j+s://` protocol for secure SSL connections, not `bolt://`

## Features

- **Connection Management**: Automatic connection pooling and session management
- **Constraints**: Unique constraints and indexes for performance
- **Error Handling**: Comprehensive error handling with context
- **Type Safety**: Strong typing with GraphQL models
- **Transaction Support**: Read and write transaction support
- **Repository Pattern**: Clean separation of concerns

## Performance Considerations

- All unique identifiers are indexed
- Common query fields (title, name, email) are indexed
- Composite constraints for rating uniqueness
- Automatic query optimization through Neo4j's query planner

## Future Enhancements

- Connection pooling configuration
- Query caching
- Database migrations
- Backup and restore functionality
- Performance monitoring
- Read replicas support
