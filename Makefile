.PHONY: help install start stop clean build-backend build-frontend dev logs

# Default target
help:
	@echo "NQ Project - Available Commands:"
	@echo ""
	@echo "  make install     - Install all dependencies (backend + frontend)"
	@echo "  make start       - Start the entire stack (backend + frontend + database)"
	@echo "  make stop        - Stop all running services"
	@echo "  make dev         - Start development mode (backend + frontend)"
	@echo "  make logs        - Show logs from all services"
	@echo "  make clean       - Clean up all build artifacts and dependencies"
	@echo "  make build       - Build both backend and frontend"
	@echo ""
	@echo "  make backend     - Start only the backend API server"
	@echo "  make frontend    - Start only the frontend development server"
	@echo "  make database    - Initialize and seed the database"
	@echo ""

# Install all dependencies
install: install-backend install-frontend
	@echo "✅ All dependencies installed successfully!"

install-backend:
	@echo "📦 Installing backend dependencies..."
	@cd backend && cargo fetch

install-frontend:
	@echo "📦 Installing frontend dependencies..."
	@cd frontend && npm install

# Start the entire stack
start: database backend frontend
	@echo "🚀 NQ stack is running!"
	@echo "   Backend API: http://localhost:8080"
	@echo "   Frontend:    http://localhost:19006"
	@echo "   Database:    SQLite (data/app.db)"
	@echo ""
	@echo "Press Ctrl+C to stop all services"

# Development mode (with hot reload)
dev: database
	@echo "🔥 Starting development mode..."
	@echo "   Backend API: http://localhost:8080"
	@echo "   Frontend:    http://localhost:19006"
	@echo ""
	@echo "Press Ctrl+C to stop all services"
	@trap 'kill 0' SIGINT; \
	cd backend && cargo run & \
	cd frontend && npm start & \
	wait

# Start backend only
backend:
	@echo "🔧 Starting backend API server..."
	@cd backend && cargo run

# Start frontend only
frontend:
	@echo "📱 Starting frontend development server..."
	@cd frontend && npm start

# Initialize and seed database
database:
	@echo "🗄️  Initializing database..."
	@cd backend && cargo run --bin init-db 2>/dev/null || \
	(echo "Database initialization failed, but continuing..." && true)

# Build both backend and frontend
build: build-backend build-frontend
	@echo "✅ Build completed!"

build-backend:
	@echo "🔧 Building backend..."
	@cd backend && cargo build --release

build-frontend:
	@echo "📱 Building frontend..."
	@cd frontend && npm run build

# Stop all services
stop:
	@echo "🛑 Stopping all services..."
	@pkill -f "cargo run" 2>/dev/null || true
	@pkill -f "npm start" 2>/dev/null || true
	@pkill -f "expo start" 2>/dev/null || true
	@echo "✅ All services stopped"

# Show logs
logs:
	@echo "📋 Service logs will appear here..."
	@echo "Press Ctrl+C to stop watching logs"
	@trap 'kill 0' SIGINT; \
	cd backend && cargo run 2>&1 | sed 's/^/[BACKEND] /' & \
	cd frontend && npm start 2>&1 | sed 's/^/[FRONTEND] /' & \
	wait

# Clean up everything
clean: stop
	@echo "🧹 Cleaning up..."
	@cd backend && cargo clean
	@cd frontend && rm -rf node_modules package-lock.json
	@rm -rf backend/data
	@echo "✅ Cleanup completed!"

# Health check
health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || echo "❌ Backend not responding"
	@curl -s http://localhost:19006 | grep -q "NQ" && echo "✅ Frontend responding" || echo "❌ Frontend not responding"

# Quick setup for new developers
setup: install database
	@echo "🎉 Setup completed! Run 'make start' to start the entire stack."

