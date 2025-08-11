#!/bin/bash

# NQ Project - Stop Script
# This script stops all running NQ services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_status "Stopping NQ services..."

# Stop backend processes
if pkill -f "cargo run" 2>/dev/null; then
    print_status "Backend service stopped"
else
    print_status "No backend service was running"
fi

# Stop frontend processes
if pkill -f "npm start" 2>/dev/null; then
    print_status "Frontend service stopped"
else
    print_status "No frontend service was running"
fi

# Stop Expo processes
if pkill -f "expo start" 2>/dev/null; then
    print_status "Expo service stopped"
else
    print_status "No Expo service was running"
fi

# Stop any remaining Node.js processes on our ports
if lsof -ti:8080 | xargs kill -9 2>/dev/null; then
    print_status "Process on port 8080 stopped"
fi

if lsof -ti:19006 | xargs kill -9 2>/dev/null; then
    print_status "Process on port 19006 stopped"
fi

print_success "All NQ services have been stopped!"

# Check if any services are still running
echo ""
print_status "Checking for remaining services..."

if lsof -i:8080 >/dev/null 2>&1; then
    print_error "Port 8080 is still in use"
else
    print_success "Port 8080 is free"
fi

if lsof -i:19006 >/dev/null 2>&1; then
    print_error "Port 19006 is still in use"
else
    print_success "Port 19006 is free"
fi

echo ""
print_status "You can now run 'npm start' to start the services again"

