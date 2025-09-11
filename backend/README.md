# NQ API

GraphQL based API built with Go and [gqlgen](https://gqlgen.com/)

## Database Schema

The API uses SQLite with the following tables:

### Core Tables

- **users** - User accounts with authentication info
- **media_types** - Categories (Movie, TV Show, Book, Game, Music)
- **media_items** - Individual media items with metadata
- **creator_roles** - Roles (Director, Actor, Author, Developer, Artist)
- **creators** - People involved in media creation
- **platforms** - Streaming platforms and stores

### Join Tables

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
