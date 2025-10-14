# NQ Backend Integrations

This package provides a comprehensive framework for integrating with various third-party services to sync user media data. The integrations support multiple media types including music, games, videos, books, and articles.

## Overview

The integrations package is designed with a modular architecture that allows easy addition of new services while maintaining consistency across all integrations.

### Supported Services

- **Spotify** - Music streaming service (albums, playlists)
- **Steam** - Gaming platform (game library, achievements)
- **YouTube** - Video platform (videos, channels, playlists)
- **YouTube Music** - Music streaming via YouTube (music videos, playlists)
- **Twitch** - Live streaming platform (games, channels)
- **Apple Music** - Apple's music streaming service (albums, artists)
- **Instapaper** - Read-later service (articles, bookmarks)

### Media Types

Each integration supports one or more media types:

- `MediaTypeGame` - Video games
- `MediaTypeMusic` - Music albums and tracks
- `MediaTypeVideo` - Videos and movies
- `MediaTypeBook` - Books and publications
- `MediaTypeArticle` - Articles and written content

## Architecture

### Core Components

1. **Integration Interface** - Defines the contract all integrations must implement
2. **Manager** - Handles registration and orchestration of multiple integrations
3. **BaseIntegration** - Provides common functionality for all integrations
4. **Specific Integrations** - Service-specific implementations

### Key Interfaces

```go
type Integration interface {
    GetName() string
    Authenticate(ctx context.Context, credentials map[string]string) error
    SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error)
    IsAuthenticated() bool
    GetSupportedMediaTypes() []MediaType
}
```

## Usage

### Basic Setup

```go
import "github.com/grillinr/nq/integrations"

// Create manager
manager := integrations.NewManager()

// Setup individual integrations
spotify := integrations.NewSpotifyIntegration()
steam := integrations.NewSteamIntegration()

// Authenticate
ctx := context.Background()
err := spotify.Authenticate(ctx, map[string]string{
    "client_id": "your-spotify-client-id",
    "client_secret": "your-spotify-client-secret",
})

// Register with manager
manager.RegisterIntegration(spotify)
manager.RegisterIntegration(steam)
```

### Syncing User Data

```go
userID := uuid.New()

// Sync from all authenticated integrations
results, err := manager.SyncAllUserData(ctx, userID)
if err != nil {
    log.Printf("Sync errors: %v", err)
}

// Process results
for integrationName, result := range results {
    fmt.Printf("Synced %d items from %s\n", 
        result.ItemsProcessed, integrationName)
    
    // Access media data
    for mediaType, items := range result.MediaData {
        fmt.Printf("Found %d %s items\n", len(items), mediaType)
    }
}
```

### Individual Integration Usage

```go
// Spotify example
spotify := integrations.NewSpotifyIntegration()
err := spotify.Authenticate(ctx, nil) // Uses env vars
if err == nil {
    result, err := spotify.SyncUserData(ctx, userID)
    // Process result...
}

// Steam example (requires Steam ID)
steam := integrations.NewSteamIntegration()
credentials := map[string]string{
    "steam_id": "your-steam-id",
}
err = steam.Authenticate(ctx, credentials)
```

## Configuration

### Environment Variables

Each integration requires specific environment variables for authentication:

#### Spotify
```bash
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
```

#### Steam
```bash
STEAM_API_KEY=your_api_key
```

#### YouTube/YouTube Music
```bash
YOUTUBE_API_KEY=your_api_key
```

#### Twitch
```bash
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
```

#### Apple Music
```bash
APPLE_MUSIC_DEVELOPER_TOKEN=your_jwt_token
```

#### Instapaper
```bash
INSTAPAPER_USERNAME=your_username
INSTAPAPER_PASSWORD=your_password
```

### Authentication Methods

Different services use different authentication methods:

- **Client Credentials** (Spotify, Twitch) - App-level access using client ID/secret
- **API Key** (Steam, YouTube) - Simple API key authentication
- **JWT Token** (Apple Music) - Developer token for catalog access
- **Basic Auth** (Instapaper) - Username/password authentication

## Data Flow

1. **Authentication** - Each integration authenticates with its respective service
2. **Data Fetching** - Service-specific APIs are called to retrieve user data
3. **Normalization** - Raw API responses are converted to standardized formats
4. **Processing** - Data is structured into common media item formats
5. **Storage** - Processed data can be stored in your application's database

## Extending the System

### Adding a New Integration

1. **Create the integration struct**:
```go
type NewServiceIntegration struct {
    *BaseIntegration
    client *http.Client
    // service-specific fields
}
```

2. **Implement the Integration interface**:
```go
func (n *NewServiceIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
    // Authentication logic
}

func (n *NewServiceIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
    // Data syncing logic
}
```

3. **Add conversion functions**:
```go
func (n *NewServiceIntegration) convertToMediaMap(item ServiceItem) map[string]interface{} {
    // Convert service data to standard format
}
```

### Supported Data Fields

The system supports these common fields across all media types:

- **Common**: `title`, `description`, `external_id`, `source`, `url`, `cover_url`, `release_date`
- **Games**: `genres`, `developers`, `publishers`, `esrb_rating`, `multiplayer`
- **Music**: `artist`, `album`, `track_count`, `duration`, `label`, `genres`
- **Videos**: `channel`, `duration`, `view_count`, `like_count`
- **Books**: `isbn`, `pages`, `publisher`, `authors`
- **Articles**: `reading_status`, `progress`, `starred`, `added_at`

## Error Handling

The system provides comprehensive error handling:

- **Authentication errors** - Invalid credentials, expired tokens
- **API errors** - Rate limiting, service unavailable
- **Data errors** - Invalid responses, missing required fields
- **Network errors** - Timeouts, connection issues

Errors are collected and returned in the `SyncResult.Errors` field, allowing partial success scenarios.

## Testing

Run the test suite:

```bash
go test ./integrations -v
```

The tests cover:
- Integration manager functionality
- Base integration behavior
- Individual integration initialization
- Data conversion utilities
- Error handling scenarios

## Rate Limiting

Each integration handles rate limiting according to the service's requirements:

- **Spotify**: Uses retry logic with exponential backoff
- **Steam**: Respects API rate limits with delays
- **YouTube**: Implements quota management
- **Others**: Service-specific rate limiting strategies

## Security Considerations

- **Credentials**: Never log or expose API keys/secrets
- **Tokens**: Implement proper token refresh mechanisms
- **Data**: Sanitize and validate all external data
- **Access**: Use least-privilege API scopes when available

## Future Enhancements

Planned improvements include:

- **OAuth 2.0 Flow**: Full OAuth support for user-specific data access
- **Webhook Support**: Real-time data updates via webhooks
- **Caching**: Intelligent caching to reduce API calls
- **Batch Processing**: Bulk operations for large datasets
- **Analytics**: Integration usage and performance metrics
- **More Services**: Netflix, Amazon Prime, Goodreads, etc.

## Troubleshooting

### Common Issues

1. **Authentication failures**: Check API keys and credentials
2. **Empty results**: Verify user has data in the service
3. **Rate limiting**: Implement delays between requests
4. **Token expiry**: Ensure token refresh mechanisms work

### Debug Mode

Enable debug logging:
```go
// Add debug logging to your integration
log.Printf("Fetching data from %s for user %s", integration.GetName(), userID)
```

## Contributing

When adding new integrations:

1. Follow the existing patterns and interfaces
2. Include comprehensive tests
3. Add proper error handling
4. Document the required environment variables
5. Update this README with the new service information