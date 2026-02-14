# NQ Wiki

Welcome to the NQ project wiki! This documentation will help you install, configure, and use NQ.

## About NQ

NQ (Next in Queue) is a smart media recommendation system that helps you decide what to watch, play, read, or listen to next. Instead of endlessly scrolling through streaming services or your game library, NQ learns your preferences and suggests content you'll actually enjoy.

**How it works:**
1. Connect your media accounts (Spotify, Steam, YouTube, etc.)
2. NQ syncs your viewing/playing/listening history
3. A graph database maps relationships between media items
4. Get personalized recommendations based on your actual preferences

## Quick Start

**New to NQ?** Start here:

1. **[User Guide](User-Guide.md)** - Complete beginner-friendly walkthrough
2. **[Installation Guide](Installation.md)** - Step-by-step setup instructions
3. **[FAQ](FAQ.md)** - Common questions and troubleshooting

**Ready to customize?** Check these out:

4. **[API Data Sources](API-Data-Sources.md)** - Connect Spotify, Steam, YouTube, and more
5. **[Tech Stack](Tech-Stack.md)** - Technologies and architecture
6. **[Usage Guide](Usage.md)** - Development workflow and advanced features

## Documentation Overview

### For End Users

- **[User Guide](User-Guide.md)** - Everything you need to know to use NQ
  - Installation from scratch
  - Connecting your media accounts
  - Getting recommendations
  - Troubleshooting
  
- **[Installation Guide](Installation.md)** - Quick setup reference
  - Installing prerequisites (Go, Node.js, Neo4j)
  - Setting up backend and frontend
  - Verification steps

- **[FAQ](FAQ.md)** - Frequently Asked Questions
  - General questions about NQ
  - Installation help
  - API and configuration questions
  - Troubleshooting common issues
  - Privacy and data questions

### For Developers

- **[API Data Sources](API-Data-Sources.md)** - External API integration reference
  - Detailed documentation for all 10 integrated APIs
  - How to get credentials for each service
  - Authentication methods
  - Rate limits and best practices

- **[Tech Stack](Tech-Stack.md)** - Technical architecture
  - Backend technologies (Go, GraphQL, Neo4j)
  - Frontend technologies (React Native, Expo, TypeScript)
  - Database schema
  - Development tools

- **[Usage Guide](Usage.md)** - Development and usage
  - Running the application
  - Development workflow
  - Testing procedures
  - Common development tasks

## Project Overview

- **Backend**: Go-based GraphQL API with Neo4j graph database
- **Frontend**: React Native mobile app built with Expo
- **Goal**: Simplify media discovery through intelligent, personalized recommendations

## Key Features

- **Unified Tracking** - Aggregate media consumption across multiple platforms
- **Smart Recommendations** - Graph-based engine analyzes relationships and preferences
- **10+ Integrations** - Connect Spotify, Steam, YouTube, Twitch, TMDB, and more
- **Cross-Platform** - React Native app works on iOS and Android
- **GraphQL API** - Flexible, self-documenting API for custom integrations
- **Privacy-First** - Self-hosted, your data stays on your infrastructure

## Media Coverage

NQ supports tracking and recommendations for:

- **Music** - Spotify, Apple Music, YouTube Music
- **Games** - Steam library, IGDB metadata, Twitch trends
- **Video** - YouTube, Movies & TV via TMDB
- **Books** - Open Library metadata and ISBN lookups
- **Articles** - Instapaper reading lists

See [API Data Sources](API-Data-Sources.md) for complete details on each integration.

## Project Structure

```
nq/
├── backend/          # Go GraphQL API server
│   ├── integrations/ # Third-party service integrations (Spotify, Steam, etc.)
│   ├── metadata/     # Media metadata providers (TMDB, IGDB, Open Library)
│   ├── db/           # Database models and repositories
│   └── graph/        # GraphQL schema and generated code
├── nq-frontend/      # React Native mobile app (Expo)
└── wiki/             # Documentation (you are here!)
```

## Getting Started

### I'm a user who wants to use NQ

👉 **Start with the [User Guide](User-Guide.md)**

This complete walkthrough covers:
- Installing all required software
- Setting up your database
- Configuring NQ
- Connecting your media accounts
- Getting your first recommendations

### I'm a developer who wants to contribute

👉 **Start with the [Installation Guide](Installation.md), then read [Usage Guide](Usage.md)**

Learn about:
- Development workflow
- Running tests
- GraphQL API structure
- Adding new integrations
- Code standards and practices

### I have a specific question

👉 **Check the [FAQ](FAQ.md)**

Common topics covered:
- Troubleshooting installation issues
- API key setup and configuration
- Performance and optimization
- Privacy and data security
- Feature requests and bug reports

## Getting Help

**Installation problems?** See the [Installation Guide](Installation.md) or [FAQ](FAQ.md)

**Need API credentials?** See [API Data Sources](API-Data-Sources.md)

**Want to understand the tech?** See [Tech Stack](Tech-Stack.md)

**Development questions?** See [Usage Guide](Usage.md)

**Can't find your answer?** Submit an issue on the [GitHub repository](https://github.com/grillinr/nq/issues)

## System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux
- **Go**: Version 1.25.1 or higher
- **Node.js**: Version 18 or higher
- **Database**: Neo4j Aura (cloud) or Neo4j Desktop (local)
- **Disk Space**: ~2GB (including dependencies and database)
- **Memory**: 4GB RAM recommended

## License

NQ is open source software released under the MIT License. See the LICENSE file in the repository for details.
