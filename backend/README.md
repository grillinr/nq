# NQ API

GraphQL based API built with Go and [gqlgen](https://gqlgen.com/)

## Quick Start

1. **Configure environment**: Copy `env.example` to `.env` and update with your Aura credentials
2. **Run the server**: `go run .`
3. **Access GraphQL Playground**: Visit `http://localhost:8080`

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
