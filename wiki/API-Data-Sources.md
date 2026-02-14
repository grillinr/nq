# API Data Sources

NQ integrates with 10 external APIs to aggregate media data and user activity across multiple platforms. These APIs are divided into two categories: **User Data Integrations** and **Metadata Providers**.

## Overview

- **User Data Integrations (7 APIs)**: Sync user activity, preferences, and consumption history
- **Metadata Providers (3 APIs)**: Enrich media items with detailed information and relationships

---

## User Data Integrations

These APIs track what media users are consuming, saving, or interacting with across different platforms.

### Spotify Web API

**Purpose**: Music streaming platform integration

**Data Collected**:
- User's saved albums
- Album metadata (title, artists, track count, release date)
- Cover artwork
- Record labels

**Authentication**: OAuth 2.0 Client Credentials Flow

**Environment Variables**:
```bash
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
```

**How to Get Credentials**:
1. Visit [Spotify for Developers](https://developer.spotify.com/dashboard)
2. Create a new application
3. Copy Client ID and Client Secret

**Implementation**: `/home/nathan/repos/nq/backend/integrations/spotify.go`

**Rate Limiting**: Implements retry logic with exponential backoff

---

### Steam Web API

**Purpose**: Gaming platform library tracking

**Data Collected**:
- User's owned games library
- Game titles, descriptions, genres
- Developers and publishers
- Playtime statistics
- Cover images and icons
- Release dates

**Authentication**: API Key + Steam ID (per user)

**Environment Variables**:
```bash
STEAM_API_KEY=your_api_key
```

**How to Get Credentials**:
1. Visit [Steam Web API Key](https://steamcommunity.com/dev/apikey)
2. Register for a key
3. Each user must provide their Steam ID

**Implementation**: `/home/nathan/repos/nq/backend/integrations/steam.go`

**Rate Limiting**: Implements delays between requests

---

### YouTube Data API v3

**Purpose**: Video platform content tracking

**Data Collected**:
- Video metadata (title, description, channel, duration)
- View counts, like counts, comment counts
- Thumbnails
- Publication dates

**Authentication**: API Key

**Environment Variables**:
```bash
YOUTUBE_API_KEY=your_api_key
```

**How to Get Credentials**:
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable YouTube Data API v3
4. Create API credentials

**Implementation**: `/home/nathan/repos/nq/backend/integrations/youtube.go`

**Limitations**: Currently demonstrates public video search. User-specific data (liked videos, subscriptions) requires OAuth 2.0.

**API Quota**: Subject to YouTube API daily quota limits

---

### YouTube Music API

**Purpose**: Music streaming via YouTube platform

**Data Collected**:
- Music videos and playlists
- Artist information
- Duration and thumbnails
- Tags and genres

**Authentication**: API Key (same as YouTube Data API)

**Environment Variables**:
```bash
YOUTUBE_API_KEY=your_api_key
```

**How to Get Credentials**: Same as YouTube Data API

**Implementation**: `/home/nathan/repos/nq/backend/integrations/yt_music.go`

**Technical Note**: Uses YouTube Data API v3 with music-specific category filtering (categoryId=10)

---

### Twitch API (Helix)

**Purpose**: Live streaming and gaming platform tracking

**Data Collected**:
- Top games on Twitch
- Game box art
- User followed channels (requires user OAuth)

**Authentication**: OAuth 2.0 Client Credentials Flow

**Environment Variables**:
```bash
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
```

**How to Get Credentials**:
1. Visit [Twitch Developers Console](https://dev.twitch.tv/console)
2. Register a new application
3. Copy Client ID and Client Secret

**Implementation**: `/home/nathan/repos/nq/backend/integrations/twitch.go`

**Token Management**: App access tokens last ~60 days, automatically refreshed at 55 days

**Limitations**: Many user-specific endpoints require user OAuth tokens (not implemented yet)

---

### Apple Music API

**Purpose**: Apple's music streaming service integration

**Data Collected**:
- Album catalog data
- Artist information, track counts
- Record labels, genres, copyright
- Artwork URLs
- Release dates

**Authentication**: JWT Developer Token

**Environment Variables**:
```bash
APPLE_MUSIC_DEVELOPER_TOKEN=your_jwt_token
```

**How to Get Credentials**:
1. Enroll in [Apple Developer Program](https://developer.apple.com/programs/)
2. Create a MusicKit identifier and private key
3. Generate a JWT developer token

**Implementation**: `/home/nathan/repos/nq/backend/integrations/apple_music.go`

**Limitations**: User library access requires additional MusicKit JS or user authentication

---

### Instapaper API

**Purpose**: Read-later service for articles and web content

**Data Collected**:
- User's saved bookmarks/articles
- Article titles, URLs, descriptions
- Reading progress and status
- Starred status
- Timestamps (added, last read)

**Authentication**: Basic Auth (username/password)

**Environment Variables**:
```bash
INSTAPAPER_USERNAME=your_username
INSTAPAPER_PASSWORD=your_password
```

**How to Get Credentials**: Use your existing Instapaper account credentials

**Implementation**: `/home/nathan/repos/nq/backend/integrations/instapaper.go`

**Features**: Tracks reading progress (unread, in_progress, completed)

---

## Metadata Providers

These APIs provide enriched metadata for media items discovered through user integrations or direct search.

### TMDB (The Movie Database)

**Purpose**: Comprehensive movie and TV show metadata

**Data Collected**:
- Movie/TV show details (title, description, release dates)
- Genres, budget, box office revenue, runtime
- Cast and crew credits with person IDs
- Production companies and countries
- Poster images and artwork
- Similar titles
- Person filmography

**Authentication**: API Read Access Token (Bearer token)

**Environment Variables**:
```bash
TMDB_API_READ_ACCESS_TOKEN=your_bearer_token
```

**How to Get Credentials**:
1. Create account at [TMDB](https://www.themoviedb.org/)
2. Visit [API Settings](https://www.themoviedb.org/settings/api)
3. Request an API key
4. Use the v4 API Read Access Token

**Implementation**: `/home/nathan/repos/nq/backend/metadata/video.go`

**Library Used**: `github.com/cyruzin/golang-tmdb` - Official Go client

**Features**:
- Auto-retry on rate limiting
- Language-specific searches
- Comprehensive credit data with character/role information

---

### IGDB (Internet Game Database)

**Purpose**: Comprehensive game metadata database

**Data Collected**:
- Game details (title, summary, release dates)
- Genres, themes, keywords
- Game modes, player perspectives
- Franchises, platforms
- Cover artwork
- Developers, publishers

**Authentication**: OAuth 2.0 via Twitch (IGDB is owned by Twitch)

**Environment Variables**:
```bash
IGDB_CLIENT_ID=your_twitch_client_id
IGDB_CLIENT_SECRET=your_twitch_client_secret
```

**How to Get Credentials**:
1. Visit [Twitch Developers Console](https://dev.twitch.tv/console)
2. Register a new application
3. Copy Client ID and Client Secret
4. Use these for IGDB API access

**Implementation**: `/home/nathan/repos/nq/backend/metadata/games.go`

**Features**:
- Automatic token refresh
- Request throttling (300ms minimum interval)
- 10-minute response caching
- Retry logic for rate limiting
- Custom query language support

**Technical Note**: IGDB requires Twitch credentials since Twitch acquired the service

---

### Open Library API

**Purpose**: Free, open book metadata database

**Data Collected**:
- Book details (title, authors, publishers)
- ISBN lookups (ISBN-10 and ISBN-13)
- Publication dates, page counts
- Cover images
- Subjects (topics, places, people, times)
- Language information

**Authentication**: None required (public API)

**Environment Variables**: None needed

**How to Use**: No credentials required - publicly accessible

**Implementation**: `/home/nathan/repos/nq/backend/metadata/books.go`

**Features**:
- ISBN-13 preferred over ISBN-10
- Multi-language support with ISO 639 conversion
- Smart title normalization for better search results
- Subject categorization (places, people, times)

**API Endpoints**:
- Book lookup: `https://openlibrary.org/api/books`
- Search: `https://openlibrary.org/search.json`
- Cover images: `https://covers.openlibrary.org/b/id/{id}-{size}.jpg`

---

## Configuration Guide

### Environment Setup

1. Copy the environment template:
   ```bash
   cd backend
   cp .envtemplate .env
   ```

2. Fill in credentials for the APIs you want to use:
   ```bash
   # User Data Integrations
   SPOTIFY_CLIENT_ID=your_value
   SPOTIFY_CLIENT_SECRET=your_value
   APPLE_MUSIC_DEVELOPER_TOKEN=your_value
   YOUTUBE_API_KEY=your_value
   TWITCH_CLIENT_ID=your_value
   TWITCH_CLIENT_SECRET=your_value
   STEAM_API_KEY=your_value
   INSTAPAPER_USERNAME=your_value
   INSTAPAPER_PASSWORD=your_value
   
   # Metadata Providers
   TMDB_API_READ_ACCESS_TOKEN=your_value
   IGDB_CLIENT_ID=your_value
   IGDB_CLIENT_SECRET=your_value
   # Open Library requires no credentials
   ```

3. The system gracefully handles missing credentials - only configured APIs will be active

### Authentication Methods Summary

| API | Method | Credentials Needed |
|-----|--------|-------------------|
| Spotify | OAuth 2.0 Client Credentials | Client ID + Secret |
| Steam | API Key | API Key + User Steam ID |
| YouTube | API Key | Google API Key |
| YouTube Music | API Key | Google API Key |
| Twitch | OAuth 2.0 Client Credentials | Client ID + Secret |
| Apple Music | JWT Token | Developer Token |
| Instapaper | Basic Auth | Username + Password |
| TMDB | Bearer Token | Read Access Token |
| IGDB | OAuth 2.0 (Twitch) | Twitch Client ID + Secret |
| Open Library | None | No credentials required |

---

## Data Flow Architecture

### Integration Framework

**Location**: `/home/nathan/repos/nq/backend/integrations/integrations.go`

The integration framework provides:
- Common interface for all integrations
- Manager pattern for orchestrating multiple services
- Standardized sync results with error collection
- Media type categorization (game, music, video, book, article)

### Metadata Service

**Location**: `/home/nathan/repos/nq/backend/metadata/metadata.go`

The metadata service provides:
- Unified service for all metadata fetchers
- Automatic fetcher initialization based on available credentials
- Graceful degradation when API keys are missing
- Language-aware metadata retrieval

### Error Handling

All integrations implement robust error handling:
- Partial success scenarios supported
- Missing credentials handled gracefully
- Rate limiting and retry logic
- Error collection in sync results

---

## Media Coverage by API

### Music
- **Spotify**: User's saved albums
- **Apple Music**: Catalog search and user library
- **YouTube Music**: Music videos and playlists

### Games
- **Steam**: User's owned games and playtime
- **Twitch**: Popular games and streaming content
- **IGDB**: Comprehensive game metadata

### Video
- **YouTube**: Video content and statistics
- **Twitch**: Live streaming content
- **TMDB**: Movies and TV shows

### Books
- **Open Library**: Book metadata and ISBN lookups

### Articles
- **Instapaper**: Saved articles and reading progress

---

## API Rate Limits & Best Practices

### Spotify
- Implements automatic retry with exponential backoff
- Token refresh at 55 minutes (tokens expire at 60 minutes)

### Steam
- Implements delays between requests
- Separate calls for library and detailed game info

### YouTube / YouTube Music
- Subject to daily quota limits (10,000 units default)
- Search operations cost 100 units each
- Video details cost 1 unit per request

### Twitch
- Token refresh at 55 days (tokens last ~60 days)
- Rate limits enforced by API

### IGDB
- 300ms minimum interval between requests
- 10-minute response caching
- Automatic rate limit retry logic
- 4 requests per second limit

### TMDB
- Auto-retry on rate limiting (40 requests per 10 seconds)
- Language-specific request support

### Open Library
- No authentication required
- Title normalization for better results

---

## Additional Resources

- **Integration Guide**: `/home/nathan/repos/nq/backend/integrations/README.md`
- **Metadata API Details**: `/home/nathan/repos/nq/backend/metadata/README.md`
- **Example Usage**: `/home/nathan/repos/nq/backend/examples/README.md`
- **Main Documentation**: [Installation Guide](Installation.md) | [Tech Stack](Tech-Stack.md) | [Usage Guide](Usage.md)

---

## Future API Integrations

Potential additions mentioned in documentation:
- Goodreads (books and reading history)
- Last.fm (music scrobbling)
- Letterboxd (movie tracking)
- PlayStation Network (gaming)
- Xbox Live (gaming)
