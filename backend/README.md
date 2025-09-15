# NQ API

GraphQL based API built with Go and [gqlgen](https://gqlgen.com/)

## Quick Start

1. **Set up Neo4j Aura**: Follow the [Aura Setup Guide](AURA_SETUP.md)
2. **Configure environment**: Copy `env.example` to `.env` and update with your Aura credentials
3. **Run the server**: `go run .`
4. **Access GraphQL Playground**: Visit `http://localhost:8080`

## Database

The API uses **Neo4j Aura** (cloud) or local Neo4j with the following structure:

### Nodes

- **users** - User accounts with authentication info
- **media_types** - Categories (Movie, TV Show, Book, Game, Music)
- **media_items** - Individual media items with metadata
- **creator_roles** - Roles (Director, Actor, Author, Developer, Artist)
- **creators** - People involved in media creation
- **platforms** - Streaming platforms and stores

### Relationships

- **media_creators** - Links media to creators
- **media_platforms** - Links media to platforms
- **media_tags** - Links media to tags
- **external_ids** - Platform-specific identifiers

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
