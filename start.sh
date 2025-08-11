#!/bin/bash

# NQ Project - Start Script
# This script starts the entire NQ stack (backend + frontend + database)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a port is in use
port_in_use() {
    lsof -i :$1 >/dev/null 2>&1
}

# Function to cleanup on exit
cleanup() {
    print_status "Shutting down services..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
    print_success "All services stopped"
    exit 0
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM

# Check prerequisites
print_status "Checking prerequisites..."

if ! command_exists cargo; then
    print_error "Rust/Cargo is not installed. Please install Rust first."
    exit 1
fi

if ! command_exists npm; then
    print_error "Node.js/npm is not installed. Please install Node.js first."
    exit 1
fi

if ! command_exists sqlite3; then
    print_warning "SQLite3 is not installed. Installing..."
    if command_exists pacman; then
        sudo pacman -S sqlite
    elif command_exists apt; then
        sudo apt update && sudo apt install -y sqlite3
    elif command_exists yum; then
        sudo yum install -y sqlite
    elif command_exists brew; then
        brew install sqlite
    else
        print_error "Could not install SQLite3 automatically. Please install it manually."
        exit 1
    fi
fi

print_success "Prerequisites check completed"

# Check if ports are available
print_status "Checking port availability..."

if port_in_use 8080; then
    print_error "Port 8080 is already in use. Please free it up first."
    exit 1
fi

if port_in_use 19006; then
    print_error "Port 19006 is already in use. Please free it up first."
    exit 1
fi

print_success "Ports are available"

# Install dependencies if needed
print_status "Installing dependencies..."

if [ ! -d "backend/target" ]; then
    print_status "Installing backend dependencies..."
    cd backend && cargo fetch && cd ..
fi

if [ ! -d "frontend/node_modules" ]; then
    print_status "Installing frontend dependencies..."
    cd frontend && npm install && cd ..
fi

print_success "Dependencies are ready"

# Initialize database
print_status "Initializing database..."
cd backend
if [ ! -f "data/app.db" ]; then
    mkdir -p data
    sqlite3 data/app.db < data_def.sql
    print_success "Database initialized"
else
    print_status "Database already exists"
fi
cd ..

# Start backend
print_status "Starting backend API server..."
cd backend
cargo run &
BACKEND_PID=$!
cd ..

# Wait for backend to start
sleep 3

# Check if backend is running
if ! curl -s http://localhost:8080/health >/dev/null; then
    print_error "Backend failed to start"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
fi

print_success "Backend started on http://localhost:8080"

# Start frontend
print_status "Starting frontend development server..."
cd frontend
npm start &
FRONTEND_PID=$!
cd ..

# Wait for frontend to start
sleep 5

# Check if frontend is running
if ! curl -s http://localhost:19006 >/dev/null; then
    print_warning "Frontend might still be starting up..."
fi

print_success "Frontend started on http://localhost:19006"

# Display status
echo ""
print_success "🎉 NQ stack is now running!"
echo ""
echo "   📍 Backend API:  http://localhost:8080"
echo "   📱 Frontend:     http://localhost:19006"
echo "   🗄️  Database:     SQLite (backend/data/app.db)"
echo ""
echo "   📋 Health Check: http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for user to stop
wait

