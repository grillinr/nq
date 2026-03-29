# nq

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white) ![React Native](https://img.shields.io/badge/react_native-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB) ![Expo](https://img.shields.io/badge/expo-1C1E24?style=for-the-badge&logo=expo&logoColor=#D04A37) ![GraphQL](https://img.shields.io/badge/-GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white) ![Neo4j](https://img.shields.io/badge/Neo4j-008CC1?style=for-the-badge&logo=neo4j&logoColor=white) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## University of Senior Design — Final Design Report

### People

Developers: Nathan Grilliot, Computer Science
Advisor: Dr. Will Hawkins, Computer Science

---

## Table of Contents

1. [Project Description](#project-description)
2. [User Interface Specification](#user-interface-specification)
3. [Test Plan and Results](#test-plan-and-results)
4. [User Manual](#user-manual)
5. [Spring Final Presentation](#spring-final-presentation)
6. [Expo Poster](#expo-poster)
7. [Assessments](#assessments)
8. [Summary of Hours and Justification](#summary-of-hours-and-justification)
9. [Summary of Expenses](#summary-of-expenses)
10. [Appendix](#appendix)

---

### 1. Project Description

**Abstract (400 characters):** NQ is a cross-platform media recommendation system that unifies user preferences across movies, TV shows, games, books, and music. It aggregates data from multiple external APIs into a Neo4j graph database, enabling personalized "next in queue" suggestions based on your consumption history. Built with Go, GraphQL, and React Native.

#### Full Project Description

NQ aims to develop a cross-platform media recommendation system that unifies user preferences across movies, shows, games, books, and more to deliver personalized "next in queue" suggestions. By aggregating and restructuring data from multiple APIs, the system will create a graph of your media consumption history. This enables context-aware recommendations based on your preferences. NQ focuses on backend functionality first, creating an API that assigns fitness scores to potential recommendations. The ultimate goal is to produce a functional tool that simplifies media discovery and adds everyday value.

- [User Stories](./sr_design_deliverables/User_Stories.md)
- [Design Diagrams](./sr_design_deliverables/Design_Diagrams/Design_Diagram.pdf)
- [Tasklist](./sr_design_deliverables/Tasklist.md)
- [Timeline and Effort Matrix](./sr_design_deliverables/Timeline.pdf)

---

### 2. User Interface Specification

- [UI Specification](./sr_design_deliverables/ui_specification.md)
- [User Guide (with screenshots and flows)](./wiki/User-Guide.md)

---

### 3. Test Plan and Results

- [Test Plan](./sr_design_deliverables/test_plan.md)
- [Frontend Testing Details](./frontend/TESTING.md)

---

### 4. User Manual

- [User Guide](./wiki/User-Guide.md)
- [Installation Guide](./wiki/Installation.md)
- [Usage Guide](./wiki/Usage.md)
- [FAQ](./wiki/FAQ.md)
- [API and Data Sources](./wiki/API-Data-Sources.md)
- [Tech Stack](./wiki/Tech-Stack.md)

---

### 5. Spring Final Presentation

- [Fall Presentation Slides](./sr_design_deliverables/Fall_Presentation.pptx)
- [Spring Presentation](./sr_design_deliverables/CS5002_Spring_Presentation.pptx)

---

### 6. Expo Poster

- [Expo Poster](./sr_design_deliverables/NQ_Senior_Design_Poster.pdf)

---

### 7. Assessments

#### Initial Self-Assessments (Fall Semester - CS5001)

- [Self Assessment Essay](./sr_design_deliverables/self_assessment_essay.docx)

#### Final Self-Assessments (Spring Semester - CS5002)

- [Final Self-Assessment](./sr_design_deliverables/CS5002_self_assessment_grillinr.docx)

---

### 8. Summary of Hours and Justification

- [Summary of Hours](./sr_design_deliverables/summary_of_hours.md)
- [Timeline and Effort Matrix](./sr_design_deliverables/Timeline.pdf)

---

### 9. Summary of Expenses

- [Summary of Expenses](./sr_design_deliverables/summary_of_expenses.md)

No costs were incurred during the development of this project.

---

### 10. Appendix

- [Full Appendix](./sr_design_deliverables/appendix.md)
- [Commit History](https://github.com/grillinr/nq/commits/main/)
- [ABET Concerns Essay](./sr_design_deliverables/abet_concerns_essay.pdf)
- [Team Contract](./sr_design_deliverables/CS5001_Team_Contract.docx)
- [Professional Biography](./sr_design_deliverables/grillinr_bio.md)
- [User Stories](./sr_design_deliverables/User_Stories.md)
- [Design Diagrams](./sr_design_deliverables/Design_Diagrams/Design_Diagram.pdf)

---

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
- `frontend/` - React Native mobile app built with Expo
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

### Neo4j Docker

A `docker-compose.yml` in the `backend/` directory runs a local Neo4j instance (browser UI on port 7474, Bolt on port 7687, credentials `neo4j/testpass`).

**Start:**

```bash
docker compose -f backend/docker-compose.yml up -d
```

**Stop (keep data):**

```bash
docker compose -f backend/docker-compose.yml stop
```

**Stop and delete volumes:**

```bash
docker compose -f backend/docker-compose.yml down -v
```

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
   cd frontend
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
