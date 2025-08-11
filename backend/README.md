# Cross-Media Tracking Platform API

A Rust-based REST API built with Actix-web for managing users, media items, ratings, favorites, and recommendations.

## Features

- **User Management**: CRUD operations for users with authentication provider support
- **Media Management**: Track movies, TV shows, books, games, and music
- **Rating System**: Rate media items with scores and timestamps
- **Activity Tracking**: Track what you're watching/reading/playing with status updates
- **Creator Management**: Track directors, actors, authors, developers, and artists
- **Platform Integration**: Link media to streaming platforms and stores
- **Tagging System**: Categorize media with custom tags
- **Recommendations**: Get personalized media recommendations
- **RESTful API**: Follows REST principles with proper HTTP status codes
- **SQLite Database**: Persistent storage with proper relationships
- **JSON Serialization**: Uses serde for request/response handling
- **Async/Await**: Built with modern Rust async patterns

## Prerequisites

- Rust 1.70+ (stable channel)
- Cargo (comes with Rust)

## Installation & Running

1. **Navigate to the backend directory:**
   ```bash
   cd backend
   ```

2. **Build the project:**
   ```bash
   cargo build
   ```

3. **Run the server:**
   ```bash
   cargo run
   ```

The API will start on `http://127.0.0.1:8080`

## Database Schema

The API uses SQLite with the following tables:

### Core Tables
- **users** - User accounts with authentication info
- **media_types** - Categories (Movie, TV Show, Book, Game, Music)
- **media_items** - Individual media items with metadata
- **creator_roles** - Roles (Director, Actor, Author, Developer, Artist)
- **creators** - People involved in media creation
- **platforms** - Streaming platforms and stores

### Junction Tables
- **media_creators** - Links media to creators
- **media_platforms** - Links media to platforms
- **media_tags** - Links media to tags
- **external_ids** - Platform-specific identifiers

### Activity & Social Tables
- **activity_statuses** - Status options (Want to Watch, Currently Watching, Completed, etc.)
- **user_activities** - User's media consumption tracking
- **ratings** - User ratings for media items
- **favorites** - User's favorite media
- **recommendations** - Media recommendations
- **tags** - Custom categorization tags

## API Endpoints

### Health Check
- **GET** `/health` - Check API status and database info

### Users
- **GET** `/users` - List all users
- **POST** `/users` - Create a new user
- **GET** `/users/{id}` - Get user by ID
- **PATCH** `/users/{id}` - Update user
- **DELETE** `/users/{id}` - Delete user

### Media Items
- **GET** `/media` - List all media items
- **POST** `/media` - Create a new media item
- **GET** `/media/{id}` - Get media item by ID

### Ratings
- **POST** `/users/{id}/ratings` - Create/update a user rating

### User Activities
- **GET** `/users/{id}/activities` - Get user's media activities
- **POST** `/users/{id}/activities` - Create a new user activity

## Example Usage

### Create a User
```bash
curl -X POST http://127.0.0.1:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "auth_provider": "google"
  }'
```

### Create a Media Item
```bash
curl -X POST http://127.0.0.1:8080/media \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "type_id": 1,
    "release_date": "1999-03-31",
    "description": "A computer hacker learns about the true nature of reality.",
    "cover_url": "https://example.com/matrix.jpg"
  }'
```

### Rate a Media Item
```bash
curl -X POST "http://127.0.0.1:8080/users/USER_ID/ratings" \
  -H "Content-Type: application/json" \
  -d '{
    "media_id": "MEDIA_ID",
    "score": 9.5
  }'
```

### Track Media Activity
```bash
curl -X POST "http://127.0.0.1:8080/users/USER_ID/activities" \
  -H "Content-Type: application/json" \
  -d '{
    "media_id": "MEDIA_ID",
    "status_id": 3,
    "rating": 9.5,
    "review": "Amazing movie!",
    "started_at": "2025-08-10",
    "finished_at": "2025-08-11"
  }'
```

### Get User Activities
```bash
curl http://127.0.0.1:8080/users/USER_ID/activities
```

## Data Models

### User
```json
{
  "user_id": "uuid-string",
  "name": "string",
  "email": "string",
  "auth_provider": "string (optional)"
}
```

### Media Item
```json
{
  "media_id": "uuid-string",
  "title": "string",
  "type_id": "integer",
  "release_date": "string (optional)",
  "description": "string (optional)",
  "cover_url": "string (optional)"
}
```

### Rating
```json
{
  "user_id": "uuid-string",
  "media_id": "uuid-string",
  "score": "number",
  "rated_at": "timestamp"
}
```

### User Activity
```json
{
  "activity_id": "uuid-string",
  "user_id": "uuid-string",
  "media_id": "uuid-string",
  "status_id": "integer",
  "rating": "number (optional)",
  "review": "string (optional)",
  "started_at": "string (optional)",
  "finished_at": "string (optional)",
  "source_platform": "uuid (optional)"
}
```

## Reference Data

The API automatically seeds the database with:

### Media Types
- Movie, TV Show, Book, Game, Music

### Creator Roles
- Director, Actor, Author, Developer, Artist

### Activity Statuses
- Want to Watch/Read/Play
- Currently Watching/Reading/Playing
- Completed
- Dropped
- On Hold

### Platforms
- Netflix, Amazon Prime, Steam

## Development

### Project Structure
```
backend/
├── Cargo.toml          # Dependencies and project configuration
├── src/
│   ├── main.rs         # Main application code and routes
│   ├── db.rs           # Database operations and models
│   └── api_spec.yaml   # OpenAPI specification
├── data/               # SQLite database files (auto-created)
└── README.md           # This file
```

### Adding New Endpoints

1. Define the data structures with `#[derive(Serialize, Deserialize)]`
2. Create the handler function in `main.rs`
3. Add the route to the `App::new()` configuration
4. Add corresponding database methods in `db.rs`

### Database Operations

All database operations are handled through the `Database` struct in `db.rs`:
- Connection management with thread-safe Arc<Mutex<Connection>>
- Automatic table creation and schema management
- Proper foreign key constraints and relationships
- Transaction support for complex operations

## Testing

The API includes comprehensive sample data for testing:
- 3 sample users (John Doe, Jane Smith, Bob Johnson)
- Reference data for media types, creator roles, and activity statuses
- Platform information for popular streaming services

## Future Enhancements

- Database migrations for schema updates
- Authentication and authorization middleware
- Input validation and sanitization
- Rate limiting and API key management
- Advanced search and filtering
- Bulk operations for data import/export
- Real-time notifications
- Integration with external media APIs
- Docker containerization
- Unit and integration tests

## License

See the main project LICENSE file.
