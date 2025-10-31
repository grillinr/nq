# nq

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white) ![React Native](https://img.shields.io/badge/react_native-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB) ![Expo](https://img.shields.io/badge/expo-1C1E24?style=for-the-badge&logo=expo&logoColor=#D04A37) ![GraphQL](https://img.shields.io/badge/-GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white) ![Neo4j](https://img.shields.io/badge/Neo4j-008CC1?style=for-the-badge&logo=neo4j&logoColor=white) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

NQ aims to develop a cross-platform media recommendation system that unifies user preferences across movies, shows, games, books, and more to deliver personalized "next in queue" suggestions. By aggregating and restructuring data from multiple APIs, the system will create a graph of your media consumption history. This enables context-aware recommendations based on your preferences. NQ will focus on backend functionality first, creating an API that assigns fitness scores to potential recommendations. The ultimate goal is to produce a functional tool that simplifies media discovery and adds everyday value.

## Languages and Tools

- **Backend**: Go, GraphQL (gqlgen), Neo4j
- **Frontend**: TypeScript, React Native, Expo
- **Database**: Neo4j Aura or local Neo4j

## External APIs and Integrations

The system integrates with the following external APIs for media data and user activity:

- **IGDB (Internet Game Database)** - Game metadata and information
- **Open Library** - Book metadata and ISBN lookups
- **YouTube Data API** - Video content and playlists
- **YouTube Music API** - Music tracks and albums
- **Twitch API** - Streaming and gaming activity
- **Spotify Web API** - Music listening history and recommendations
- **Steam Web API** - Game library and playtime data
- **Apple Music API** - Music catalog and user data
- **Instapaper API** - Reading list and article bookmarks

## Project Structure

- `backend/` - Go GraphQL API server with database repositories and resolvers
- `nq-frontend/` - React Native mobile app built with Expo
- `db/` - Database models, constraints, and repository implementations
- `integrations/` - Third-party service integrations (Spotify, YouTube Music, Twitch, etc.)
- `metadata/` - Media metadata providers for books, games, movies, and TV shows
- `examples/` - Example code and integration demos
- `graph/` - GraphQL schema and generated code
- `sr_design_deliverables/` - Design documents and project deliverables

## Getting Started

### Prerequisites

- Go 1.19+
- Node.js 18+
- Neo4j (Aura cloud instance or local installation)

### Backend Setup

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Copy the environment template:

   ```bash
   cp .envtemplate .env
   ```

3. Update `.env` with your Neo4j credentials and other configuration.

4. Run the server:

   ```bash
   go run .
   ```

5. Access the GraphQL playground at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:

   ```bash
   cd nq-frontend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:

   ```bash
   npx expo start
   ```

4. Follow the prompts to open in emulator, simulator, or Expo Go app.

## Database Schema

The system uses a graph database with the following main entities:

- **Users** - User accounts
- **Media Items** - Movies, TV shows, books, games, music
- **Creators** - Actors, directors, authors, developers
- **Platforms** - Streaming services and stores
- **Tags** - Content categorization

Relationships connect these entities to enable personalized recommendations.
