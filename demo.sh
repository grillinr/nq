#!/bin/bash

# NQ Project - Demo Script
# This script demonstrates the entire NQ stack

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${PURPLE}║                    NQ Project Demo                           ║${NC}"
    echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_step() {
    echo -e "${CYAN}🔹 Step $1:${NC} $2"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to wait for service to be ready
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1
    
    print_info "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" >/dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to start within $max_attempts seconds"
    return 1
}

# Function to test API endpoints
test_api() {
    local base_url=$1
    
    print_step "1" "Testing API health endpoint"
    if curl -s "$base_url/health" | grep -q "healthy"; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        return 1
    fi
    
    print_step "2" "Testing users endpoint"
    if curl -s "$base_url/users" >/dev/null; then
        print_success "Users endpoint working"
    else
        print_error "Users endpoint failed"
        return 1
    fi
    
    print_step "3" "Testing media endpoint"
    if curl -s "$base_url/media" >/dev/null; then
        print_success "Media endpoint working"
    else
        print_error "Media endpoint failed"
        return 1
    fi
}

# Function to test frontend
test_frontend() {
    local url=$1
    
    print_step "1" "Testing frontend accessibility"
    if curl -s "$url" | grep -q "NQ"; then
        print_success "Frontend is accessible"
    else
        print_warning "Frontend might still be starting up"
    fi
}

# Main demo function
run_demo() {
    print_header
    
    print_info "This demo will start the entire NQ stack and test all services"
    echo ""
    
    # Check if services are already running
    if lsof -i:8080 >/dev/null 2>&1 || lsof -i:19006 >/dev/null 2>&1; then
        print_warning "Some services are already running. Stopping them first..."
        ./stop.sh
        sleep 2
    fi
    
    print_step "1" "Starting the entire NQ stack"
    print_info "This will start: Database → Backend API → Frontend"
    echo ""
    
    # Start services in background
    ./start.sh &
    START_PID=$!
    
    # Wait a bit for services to start
    sleep 10
    
    print_step "2" "Testing backend API (http://localhost:8080)"
    if wait_for_service "http://localhost:8080/health" "Backend API"; then
        test_api "http://localhost:8080"
    else
        print_error "Backend API test failed"
        kill $START_PID 2>/dev/null || true
        exit 1
    fi
    
    print_step "3" "Testing frontend (http://localhost:19006)"
    if wait_for_service "http://localhost:19006" "Frontend"; then
        test_frontend "http://localhost:19006"
    else
        print_warning "Frontend test incomplete (might still be starting)"
    fi
    
    print_step "4" "Displaying service status"
    echo ""
    print_success "🎉 NQ Stack Demo Completed Successfully!"
    echo ""
    echo -e "${GREEN}📍 Services Running:${NC}"
    echo "   🔧 Backend API:  http://localhost:8080"
    echo "   📱 Frontend:     http://localhost:19006"
    echo "   🗄️  Database:     SQLite (backend/data/app.db)"
    echo ""
    echo -e "${GREEN}📋 Test Endpoints:${NC}"
    echo "   Health Check:   http://localhost:8080/health"
    echo "   Users API:      http://localhost:8080/users"
    echo "   Media API:      http://localhost:8080/media"
    echo ""
    echo -e "${YELLOW}💡 Next Steps:${NC}"
    echo "   • Open http://localhost:19006 in your browser"
    echo "   • Test the mobile app on your device"
    echo "   • Explore the API endpoints"
    echo "   • Run 'npm run stop' to stop all services"
    echo ""
    
    # Keep the demo running
    print_info "Demo is running. Press Ctrl+C to stop all services and exit."
    wait $START_PID
}

# Cleanup function
cleanup() {
    print_info "Stopping demo..."
    ./stop.sh
    print_success "Demo stopped"
    exit 0
}

# Set trap for cleanup
trap cleanup SIGINT SIGTERM

# Check if we're in the right directory
if [ ! -f "start.sh" ] || [ ! -f "package.json" ]; then
    print_error "Please run this script from the NQ project root directory"
    exit 1
fi

# Run the demo
run_demo

