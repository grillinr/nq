# Tech Stack

This page describes the technologies and tools that power NQ. 

**For non-developers:** This page is more technical and intended for developers who want to understand how NQ works under the hood. If you're just trying to use NQ, you don't need to understand everything here - the [User Guide](User-Guide.md) and [Installation Guide](Installation.md) have everything you need.

**For developers:** NQ is built with modern technologies optimized for performance, scalability, and developer experience.

## Backend Technologies

### Language & Framework

- **Go 1.25.1+** - Primary backend language
- **GraphQL** - API query language for flexible data fetching
- **gqlgen (99designs/gqlgen)** - Go GraphQL server library and code generator

### Database

- **Neo4j** - Graph database for storing media relationships
  - Enables complex recommendation queries
  - Models relationships between users, media, creators, and platforms
  - Supports both Neo4j Aura (cloud) and local installations

### Key Backend Libraries

- `neo4j-go-driver/v5` - Official Neo4j driver for Go
- `golang-jwt/jwt/v5` - JWT authentication
- `joho/godotenv` - Environment variable management
- `google/uuid` - UUID generation
- `cyruzin/golang-tmdb` - The Movie Database API client

## Frontend Technologies

### Language & Framework

- **TypeScript** - Strongly-typed JavaScript for safer code
- **React Native** - Cross-platform mobile framework
- **Expo SDK 54** - React Native development platform
- **Expo Router** - File-based routing for React Native

### State Management & Data

- **Apollo Client** - GraphQL client for React
- **AsyncStorage** - Persistent local storage
- **React Hooks** - Modern state management patterns

### UI & Navigation

- **React Navigation** - Navigation library (bottom tabs, stack navigation)
- **Expo Vector Icons** - Icon library
- **React Native SVG** - SVG rendering
- **React Native Chart Kit** - Data visualization
- **React Native Gesture Handler** - Touch gesture handling

### Key Frontend Libraries

- `expo-auth-session` - OAuth authentication flows
- `expo-secure-store` - Secure credential storage
- `expo-linking` - Deep linking support
- `react-native-reanimated` - Smooth animations
- `graphql` - GraphQL query language support

## External API Integrations

NQ integrates with multiple third-party services to aggregate media data:

### Gaming

- **IGDB (Internet Game Database)** - Game metadata and information
- **Steam Web API** - Game library and playtime data
- **Twitch API** - Streaming and gaming activity

### Music

- **Spotify Web API** - Music listening history and recommendations
- **YouTube Music API** - Music tracks and albums
- **Apple Music API** - Music catalog and user data

### Video

- **YouTube Data API** - Video content and playlists
- **TMDB (The Movie Database)** - Movie and TV show metadata

### Books & Reading

- **Open Library** - Book metadata and ISBN lookups
- **Instapaper API** - Reading list and article bookmarks

## Database Schema

The Neo4j graph database models the following entities and their relationships:

### Core Entities

- **Users** - User accounts and authentication
- **Media Items** - Movies, TV shows, books, games, music
- **Creators** - Actors, directors, authors, developers, artists
- **Platforms** - Streaming services, stores, and distribution channels
- **Tags** - Content categorization and metadata

### Relationships

The graph structure connects these entities to enable:
- Personalized recommendations based on consumption history
- Discovery of similar content through shared attributes
- Context-aware suggestions considering user preferences
- Cross-platform media tracking

## Development Tools

### Backend

- **go fmt** - Code formatting
- **go vet** - Static analysis
- **go test** - Unit testing framework
- **net/http/httptest** - HTTP testing utilities

### Frontend

- **ESLint** - JavaScript/TypeScript linting
- **TypeScript Compiler** - Type checking
- **Expo CLI** - Development and build tooling

## Architecture

- **Monorepo Structure** - Backend and frontend in single repository
- **GraphQL API** - Single endpoint for all data operations
- **Graph Database** - Optimized for relationship-based queries
- **Mobile-First** - React Native for iOS and Android support
