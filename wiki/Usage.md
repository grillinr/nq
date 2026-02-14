# Usage Guide

This guide covers how to run, develop, and test the NQ project.

**For end users:** If you're looking for instructions on how to use NQ (not develop it), see the [User Guide](User-Guide.md) instead. This page is primarily for developers.

**For developers:** This page covers the development workflow, testing procedures, and advanced usage of NQ.

## Running the Application

### Backend Server

Navigate to the backend directory and start the server:

```bash
cd backend
go run .
```

The GraphQL API will be available at `http://localhost:8080`.

### Frontend App

Navigate to the frontend directory and start the Expo development server:

```bash
cd nq-frontend
npx expo start
```

This opens an interactive menu where you can:
- Press `i` to open iOS Simulator (Mac only)
- Press `a` to open Android Emulator
- Scan the QR code with Expo Go app on your phone

## Development Workflow

### Backend Development

#### Building

Compile the backend to check for errors:

```bash
cd backend
go build .
```

#### Running Tests

Run all tests:

```bash
go test ./...
```

Run a specific test:

```bash
go test -v -run TestName ./path/to/package
```

Example:

```bash
go test -v -run TestRecoverMiddleware .
```

#### Code Formatting & Linting

Format code:

```bash
go fmt ./...
```

Run static analysis:

```bash
go vet ./...
```

### Frontend Development

#### Type Checking

Verify TypeScript types without emitting files:

```bash
cd nq-frontend
npx tsc --noEmit
```

#### Linting

Check code quality and style:

```bash
npm run lint
```

#### Development Server Options

Start with cache cleared:

```bash
npx expo start -c
```

Start for specific platform:

```bash
npm run android  # Android only
npm run ios      # iOS only
npm run web      # Web browser
```

## Using the GraphQL API

### GraphQL Playground

1. Start the backend server (`go run .` in backend directory)
2. Open your browser to `http://localhost:8080`
3. The GraphQL playground interface will load
4. Write and test queries in the interactive editor

### Example Queries

The GraphQL playground provides:
- Auto-completion for queries and mutations
- Interactive schema documentation
- Query validation and error messages
- Real-time query execution

Explore the schema documentation in the playground to discover available queries and mutations.

## Environment Configuration

### Backend (.env)

The backend requires environment variables for:
- Neo4j connection (URI, username, password)
- External API keys (Spotify, YouTube, IGDB, etc.)
- Server configuration (port, environment mode)

Copy `.envtemplate` to `.env` and fill in your values:

```bash
cd backend
cp .envtemplate .env
```

Edit `.env` with your configuration.

### Frontend

Frontend configuration is managed through:
- Expo constants and environment variables
- Apollo Client configuration for GraphQL endpoint
- Secure storage for authentication tokens

## Testing

### Backend Tests

The backend uses Go's standard testing package:

```bash
cd backend
go test ./...
```

Tests use `net/http/httptest` for HTTP testing and mock repositories for database operations.

### Frontend Testing

Testing guidelines:
- Type safety through TypeScript compilation
- ESLint for code quality
- Manual testing via Expo development builds

## Common Tasks

### Adding a New GraphQL Query

1. Update schema in `graph/schema.graphqls`
2. Run `go generate ./...` to regenerate code
3. Implement resolver in appropriate resolver file
4. Test in GraphQL playground

### Updating Dependencies

Backend:

```bash
cd backend
go get -u ./...
go mod tidy
```

Frontend:

```bash
cd nq-frontend
npm update
```

### Database Migrations

Neo4j schema changes are managed through:
- Cypher queries for constraints and indexes
- Database repository implementations
- Initial setup scripts

## Development Best Practices

### Backend

- Follow Go formatting standards (`gofmt`)
- Handle errors explicitly (no silent failures)
- Use meaningful variable and function names
- Write tests for new functionality
- Document complex logic with comments

### Frontend

- Use TypeScript for type safety
- Follow ESLint configuration
- Use functional components with hooks
- Keep components small and focused
- Handle loading and error states

## Troubleshooting

### Backend Issues

**Server won't start:**
- Check Neo4j connection in `.env`
- Verify port 8080 is available
- Review error logs for specific issues

**GraphQL errors:**
- Regenerate GraphQL code: `go generate ./...`
- Check schema syntax
- Verify resolver implementations

### Frontend Issues

**Metro bundler errors:**
- Clear cache: `npx expo start -c`
- Delete `node_modules` and reinstall
- Check for conflicting dependencies

**Connection errors:**
- Ensure backend is running on port 8080
- Check Apollo Client configuration
- Verify network permissions (physical devices)

## Next Steps

- Review the [Tech Stack](Tech-Stack.md) to understand the technologies used
- Explore the codebase structure in the main README
- Check the GraphQL schema documentation in the playground
- Join the development workflow and start contributing
