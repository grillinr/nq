# Examples

This directory contains example applications demonstrating how to use various components of the NQ backend.

## Structure

```
examples/
├── integrations/          # Integration system example
│   └── main.go           # Shows how to use the integrations manager
├── metadata/             # Metadata fetching example  
│   └── main.go           # Shows how to fetch metadata for different media types
└── README.md             # This file
```

## Running Examples

You can run examples using the provided script or manually.

### Quick Start with Runner Script

```bash
# Interactive menu
./examples/run_examples.sh

# Run specific example
./examples/run_examples.sh metadata
./examples/run_examples.sh integrations

# Run all examples
./examples/run_examples.sh all
```

### Manual Execution

Each example is in its own directory and can be run independently.

#### Integration Example

Shows how to use the integrations system to sync data from various services (Spotify, Steam, YouTube, etc.).

```bash
cd examples/integrations
go run main.go
```

This example demonstrates:
- Setting up the integration manager
- Configuring different service integrations
- Authenticating with services (requires API keys)
- Syncing user data from multiple platforms
- Processing sync results and error handling

### Metadata Example

Shows how to fetch metadata for different types of media (movies, books, TV shows, games).

```bash
cd metadata
go run main.go
```

This example demonstrates:
- Fetching movie metadata by title and year
- Fetching book metadata by ISBN
- Fetching TV show metadata by title
- Using map-based metadata fetching
- Error handling for invalid or missing data

## Prerequisites

### For Integration Example
You'll need API keys for the services you want to test:
- **Spotify**: Set `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET`
- **Steam**: Set `STEAM_API_KEY`
- **YouTube**: Set `YOUTUBE_API_KEY`
- **Twitch**: Set `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET`
- **Apple Music**: Set `APPLE_MUSIC_KEY_ID`, `APPLE_MUSIC_TEAM_ID`, and `APPLE_MUSIC_PRIVATE_KEY`

### For Metadata Example
The metadata system uses public APIs:
- **TMDB** for movies and TV shows (requires `TMDB_API_KEY`)
- **Open Library** for books (no API key required)
- **IGDB** for games (requires `IGDB_CLIENT_ID` and `IGDB_CLIENT_SECRET`)

## Environment Variables

Create a `.env` file in the backend root directory:

```bash
# Metadata APIs
TMDB_API_KEY=your_tmdb_api_key
IGDB_CLIENT_ID=your_igdb_client_id
IGDB_CLIENT_SECRET=your_igdb_client_secret

# Integration APIs
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
STEAM_API_KEY=your_steam_api_key
YOUTUBE_API_KEY=your_youtube_api_key
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret
APPLE_MUSIC_KEY_ID=your_apple_music_key_id
APPLE_MUSIC_TEAM_ID=your_apple_music_team_id
APPLE_MUSIC_PRIVATE_KEY=your_apple_music_private_key_path
```

## Building Examples

You can build standalone executables for each example:

```bash
# Build integrations example
cd integrations
go build -o integrations_example main.go
./integrations_example

# Build metadata example
cd metadata
go build -o metadata_example main.go
./metadata_example
```

## Notes

- Examples may require network access to external APIs
- Some integrations require OAuth authentication flows
- Error messages will be displayed if API keys are missing or invalid
- The examples are designed to be educational and may not include production-level error handling

## Adding New Examples

To add a new example:

1. Create a new directory under `examples/`
2. Add a `main.go` file with your example code
3. Update this README to document the new example
4. Ensure your example can be built and run independently

Example structure:
```go
package main

import (
    "fmt"
    "github.com/grillinr/nq/your_package"
)

func main() {
    // Your example code here
    fmt.Println("Hello from new example!")
}
```