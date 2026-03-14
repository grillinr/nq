# Frequently Asked Questions (FAQ)

This page answers common questions about NQ installation, setup, usage, and troubleshooting.

## General Questions

### What is NQ?

NQ (Next in Queue) is a smart media recommendation system that helps you decide what to watch, play, read, or listen to next. It connects to your accounts on services like Spotify, Steam, YouTube, and more to understand your preferences and suggest content you'll enjoy.

### Who is NQ for?

NQ is designed for anyone who:
- Has trouble deciding what to watch, play, or read next
- Uses multiple streaming and media services
- Wants personalized recommendations based on their actual preferences
- Enjoys discovering new content similar to what they already like

### What platforms does NQ support?

NQ currently supports:
- **Music**: Spotify, Apple Music, YouTube Music
- **Games**: Steam, Twitch
- **Video**: YouTube, Movies/TV via TMDB
- **Books**: Open Library
- **Articles**: Instapaper

### Is NQ free to use?

Yes, NQ is free and open source under the MIT license. However, some of the external services it connects to may require their own subscriptions (like Spotify Premium or Steam).

### Do I need all the API keys to use NQ?

No! NQ is designed to work with whatever services you want to connect. You only need API keys for the services you plan to use. The system gracefully handles missing credentials.

---

## Installation Questions

### What do I need to install before using NQ?

You need three main things:
1. **Go** (version 1.25.1 or newer) - For running the backend server
2. **Node.js** (version 18 or newer) - For running the mobile app
3. **Neo4j** database - Either a free cloud instance or local installation

See the [Installation Guide](Installation.md) for detailed setup instructions.

### Which Neo4j option should I choose - cloud or local?

**Cloud (Neo4j Aura)** - Recommended for beginners:
- Free tier available
- No installation needed
- Automatic backups
- Easy setup

**Local Installation** - For developers or advanced users:
- Full control over data
- No internet required after setup
- Useful for development

### I'm not a developer. Is NQ too technical for me?

NQ is currently in development and does require some technical setup (installing software, configuring API keys). If you're comfortable following step-by-step instructions and have used command-line tools before, you should be able to get it running. Future versions will aim to simplify the setup process.

### How long does installation take?

Typical installation time:
- Installing prerequisites (Go, Node.js, Neo4j): 15-30 minutes
- Setting up the backend: 10-15 minutes
- Setting up the frontend: 5-10 minutes
- Getting API keys: Varies (5-60 minutes depending on how many services)

Total: 30 minutes to 2 hours depending on your experience level and how many integrations you want.

---

## API and Configuration Questions

### Where do I get API keys?

Each service has its own process for getting API keys. See the [API Data Sources](API-Data-Sources.md) page for detailed instructions for each service. Most are free for personal use.

### Do I have to pay for API access?

Most APIs NQ uses offer free tiers that are sufficient for personal use:
- **Free**: Open Library, YouTube, Spotify, Steam, IGDB, TMDB
- **Requires Developer Account**: Apple Music (requires Apple Developer Program membership - $99/year)
- **Free with Account**: Instapaper (uses your regular account credentials)

### What information do I put in the .env file?

The `.env` file stores your configuration and API credentials. It's like a secure settings file. Copy `.envtemplate` to `.env` and fill in your information. See the [Installation Guide](Installation.md#backend-setup) for detailed instructions.

### Is my data secure? Where are my API keys stored?

Your API keys are stored locally in the `.env` file on your computer. This file should never be shared or committed to version control (it's in `.gitignore` by default). Your data stays on your Neo4j database - either on your computer (local) or in your private cloud instance (Aura).

### I don't want to connect all services. Can I use just a few?

Absolutely! NQ works with any combination of services. Start with just one or two that you use most, and add more later if you want.

---

## Usage Questions

### How do I start using NQ after installation?

1. Start the backend server (in the `backend` folder, run `go run .`)
2. Start the mobile app (in the `frontend` folder, run `npx expo start`)
3. Open the app on your phone or emulator
4. Connect your media accounts through the app

### How does NQ make recommendations?

NQ builds a graph of your media consumption:
1. It syncs your activity from connected services (what you've watched, played, listened to)
2. It enriches this data with metadata (genres, creators, themes, etc.)
3. It analyzes relationships between media items in the graph database
4. It calculates "fitness scores" for potential recommendations based on your preferences

### How often does NQ sync my data?

This depends on your configuration and usage. You can trigger manual syncs through the app, or the system can sync periodically. The frequency is configurable.

### Can I use NQ without the mobile app?

Yes! The backend provides a GraphQL API that you can access directly through the GraphQL Playground at `http://localhost:8080`. This is useful for developers or if you want to build your own interface.

### Does NQ work offline?

The backend needs internet to sync data from external services. However, once data is synced to your Neo4j database, you can query it offline. The mobile app requires connection to your backend server.

---

## Troubleshooting

### The backend won't start

**Check these common issues:**

1. **Port 8080 already in use**
   - Another program is using port 8080
   - Solution: Stop the other program, or change NQ's port in `.env`

2. **Can't connect to Neo4j**
   - Check your Neo4j credentials in `.env`
   - Make sure your Neo4j instance is running
   - Verify the URI format: `neo4j+s://xxxxx.databases.neo4j.io` (cloud) or `bolt://localhost:7687` (local)

3. **Missing .env file**
   - You need to create it: `cp .envtemplate .env`
   - Then fill in your configuration

4. **Go not found**
   - Install Go from https://go.dev/dl/
   - Verify installation: `go version`

### The frontend won't start

**Common issues:**

1. **Packages not installed**
   - Run `npm install` in the `frontend` folder

2. **Metro bundler errors**
   - Clear cache: `npx expo start -c`
   - Delete `node_modules` and run `npm install` again

3. **Can't connect to backend**
   - Make sure the backend is running (`go run .` in backend folder)
   - Check that it's accessible at `http://localhost:8080`

4. **Node.js version too old**
   - Update to Node.js 18 or newer
   - Check version: `node --version`

### API integration isn't working

**Debugging steps:**

1. **Verify API credentials**
   - Check that the API key/credentials in `.env` are correct
   - Make sure there are no extra spaces or quotes

2. **Check API key permissions**
   - Some APIs require enabling specific permissions or scopes
   - Verify the key is active and not expired

3. **Rate limiting**
   - You may have exceeded the API's rate limit
   - Wait a bit and try again
   - Check the API's documentation for rate limits

4. **Service-specific issues**
   - Check the [API Data Sources](API-Data-Sources.md) page for service-specific troubleshooting

### GraphQL Playground shows errors

**Common causes:**

1. **Backend not running**
   - Start it: `cd backend && go run .`

2. **Wrong URL**
   - Make sure you're visiting `http://localhost:8080` (not https)

3. **Query syntax error**
   - Use the playground's auto-complete (Ctrl+Space)
   - Check the schema documentation in the playground

### Database errors

**Solutions:**

1. **"Connection refused"**
   - Neo4j isn't running
   - For Aura: Check your internet connection
   - For local: Start Neo4j Desktop or service

2. **"Authentication failed"**
   - Check username/password in `.env`
   - Neo4j Aura default username is usually `neo4j`

3. **"Database not found"**
   - Verify the database name in `.env`
   - For Aura: Usually `neo4j`
   - For local: Check Neo4j Desktop for database name

### How do I reset everything and start over?

**Complete reset:**

1. Stop the backend and frontend
2. Delete your Neo4j database (or create a new Aura instance)
3. Update `.env` with new database credentials
4. Restart the backend: `cd backend && go run .`
5. Restart the frontend: `cd frontend && npx expo start`

**Just reset API connections:**
- Remove and re-add your API credentials in `.env`
- Restart the backend

---

## Performance Questions

### Why is the initial sync taking so long?

The first sync downloads a lot of data from your connected services:
- Large libraries (1000+ games on Steam, for example) take time
- API rate limits slow down requests to prevent overwhelming services
- Metadata enrichment (fetching details for each item) adds time

This is normal! Subsequent syncs are much faster since they only update changes.

### Can I sync just one service at a time?

Yes! The integration manager allows you to sync services individually. This can be faster and easier to troubleshoot.

### How much disk space does NQ need?

- Backend + dependencies: ~100-200 MB
- Frontend + dependencies: ~500 MB - 1 GB
- Neo4j database: Varies based on library size
  - Small library (100-500 items): ~10-50 MB
  - Medium library (500-2000 items): ~50-200 MB
  - Large library (2000+ items): ~200 MB - 1 GB+

---

## Development Questions

### Can I contribute to NQ?

Yes! NQ is open source. Check out the repository and feel free to submit issues or pull requests.

### How do I run the tests?

**Backend:**
```bash
cd backend
go test ./...
```

**Frontend:**
```bash
cd frontend
npm run lint
npx tsc --noEmit
```

### Where can I find the API documentation?

The GraphQL API is self-documenting. Start the backend and visit `http://localhost:8080` to explore the interactive schema documentation.

### Can I add support for a new service?

Yes! The integration framework is designed to be extensible. See the existing integrations in `backend/integrations/` for examples. You'll need to:
1. Implement the `Integration` interface
2. Add configuration for API credentials
3. Map the service's data to NQ's media model

---

## Privacy and Data Questions

### What data does NQ collect?

NQ only collects data you explicitly authorize:
- Your media consumption history from connected services
- Public metadata about media items (titles, descriptions, etc.)
- Your preferences and interactions within the app

NQ is self-hosted, meaning all your data stays on your own infrastructure.

### Does NQ share my data with anyone?

No. NQ runs entirely on your own computer/infrastructure. It doesn't send your data to any third-party servers except when syncing from the services you've authorized (Spotify, Steam, etc.).

### Can I delete my data?

Yes. Since you control the database:
- Delete specific items through the GraphQL API
- Delete your entire database to remove all data
- Stop using NQ at any time and delete the installation

### Is my activity tracked or monitored?

No. There's no analytics, tracking, or monitoring built into NQ. It's completely private and runs locally.

---

## Future Features

### Will there be a web interface?

While NQ currently focuses on the mobile app, a web interface is possible. The GraphQL API is already built to support any client.

### Will you add support for [service]?

Potentially! Popular requests include:
- Goodreads (books)
- Last.fm (music scrobbling)
- Letterboxd (movies)
- PlayStation Network (games)
- Xbox Live (games)

Feature requests are welcome - submit them as GitHub issues.

### Is there a hosted version I can use without installing?

Not currently. NQ is designed for self-hosting to maintain privacy and control. A hosted version may be considered in the future.

---

## Getting More Help

### Where can I get help if my question isn't answered here?

1. Check the [User Guide](User-Guide.md) for step-by-step instructions
2. Review the [Installation Guide](Installation.md) for setup help
3. See [Usage Guide](Usage.md) for development workflow
4. Check the [API Data Sources](API-Data-Sources.md) for API-specific issues
5. Submit an issue on the GitHub repository

### How do I report a bug?

Submit a GitHub issue with:
- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Error messages (if any)
- Your environment (OS, Go version, Node version)

### How do I request a feature?

Submit a GitHub issue describing:
- What you'd like to add or change
- Why it would be useful
- How you envision it working

---

## Additional Resources

- [Home](Home.md) - Wiki home page
- [Installation Guide](Installation.md) - Detailed setup instructions
- [User Guide](User-Guide.md) - Step-by-step how-to guide
- [Tech Stack](Tech-Stack.md) - Technologies used
- [API Data Sources](API-Data-Sources.md) - External API documentation
- [Usage Guide](Usage.md) - Development and usage information
