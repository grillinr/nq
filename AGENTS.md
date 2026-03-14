# Agentic Development Guidelines

This document provides context and rules for AI agents operating within the `nq` repository. The project is a monorepo consisting of a Go backend and an Expo (React Native) frontend.

## 1. Project Structure

- **root**: Contains documentation.
- **`backend/`**: Go application (GraphQL API with Neo4j).
- **`frontend/`**: React Native (Expo) application.

## 2. Backend (Go)

**Location:** `/home/nathan/repos/nq/backend`

### Commands

- **Build:** `go build .`
- **Run:** `go run .` (starts server on port 8080)
- **Test (All):** `go test ./...`
- **Test (Single):** `go test -v -run TestName ./path/to/pkg` (e.g., `go test -v -run TestRecoverMiddleware .`)
- **Lint/Format:** `go fmt ./...` and `go vet ./...`

### Code Style & Conventions

- **Language:** Go 1.25.1+
- **Formatting:** Strict adherence to `gofmt`.
- **Imports:** Grouped imports (std lib first, then 3rd party, then local).
- **Naming:** CamelCase for exported members, mixedCase for private. Test functions must start with `Test`.
- **Error Handling:** Explicit error checking (`if err != nil`). Return errors to the caller unless it's `main`. Use `log` only in `main` or middleware.
- **Architecture:** GraphQL using `99designs/gqlgen`. Database interaction via `neo4j-go-driver`.
- **Testing:** Use standard `testing` package and `net/http/httptest`.

## 3. Frontend (Expo / React Native)

**Location:** `/home/nathan/repos/nq/frontend`

### Commands

- **Install:** `npm install`
- **Start:** `npm start` (Interactive Expo CLI)
- **Lint:** `npm run lint`
- **Type Check:** `npx tsc --noEmit`

### Code Style & Conventions

- **Language:** TypeScript
- **Framework:** React Native with Expo, Expo Router.
- **State/Data:** `@apollo/client` for GraphQL.
- **Formatting:** Follow ESLint config (`npm run lint` to check).
- **Components:** Functional components using Hooks. Use strictly typed props.
- **Directory Structure:**
  - `src/components/`: Reusable UI components.
  - `src/pages/`: Screen components (or strictly `app/` if using Expo Router).
  - `src/types.ts`: Shared TypeScript interfaces.
- **Naming:** PascalCase for components (`MyComponent.tsx`), camelCase for functions/vars.
- Style: Don't use inline styles. Use `StyleSheet.create` for all styles. All styling should be token-based (e.g., `color: theme.colors.primary`), no hardcoded values.

## 4. General Agent Rules

- **Path Handling:** ALWAYS use absolute paths for file operations (e.g., `/home/nathan/repos/nq/backend/main.go`).
- **Context:** Before editing, read the file and its imports to understand dependencies.
- **Verification:**
  - After backend changes: Run `go build .` and `go test ./...` in `backend/`.
  - After frontend changes: Run `npm run lint` in `frontend/`.
- **Safety:**
  - Do not commit secrets/credentials.
  - Do not revert existing code unless explicitly requested or fixing a regression you introduced.
  - If a file doesn't exist, verify the path before assuming it's missing.

## 5. Development Notes

- **Database:** The backend requires a Neo4j instance. Tests using `nil` repositories should handle the absence of a DB gracefully.
- **Port:** Backend runs on port 8080 by default.
- **Environment:** Use `.env` files for configuration (handled by `godotenv` in backend).
