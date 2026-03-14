#!/bin/bash

# NQ Development Environment Startup Script
# Starts backend and frontend in tmux session with health checks

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SESSION_NAME="nq-dev"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$REPO_ROOT/logs"
BACKEND_LOG="$LOG_DIR/backend.log"
FRONTEND_LOG="$LOG_DIR/frontend.log"
BACKEND_PORT=8080
BACKEND_HEALTH_TIMEOUT=30

# Ensure logs directory exists
mkdir -p "$LOG_DIR"

# Print header
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Starting NQ Development Environment${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${CYAN}Logs: $LOG_DIR${NC}"
echo ""

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo -e "${RED}Error: tmux is not installed${NC}"
    echo "Install it with: sudo apt install tmux (Ubuntu/Debian) or brew install tmux (macOS)"
    exit 1
fi

# Check if session already exists
if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
    echo -e "${YELLOW}Session '$SESSION_NAME' already exists${NC}"
    echo "Options:"
    echo "  1. Attach to existing session: tmux attach -t $SESSION_NAME"
    echo "  2. Kill and restart: npm run stop && npm run dev"
    exit 1
fi

# Function to wait for backend health
wait_for_backend() {
    echo -e "${YELLOW}[1/2] Starting Go backend...${NC}"
    local elapsed=0
    while [ $elapsed -lt $BACKEND_HEALTH_TIMEOUT ]; do
        if curl -fsS http://localhost:$BACKEND_PORT > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Backend healthy at http://localhost:$BACKEND_PORT${NC}"
            return 0
        fi
        sleep 1
        elapsed=$((elapsed + 1))
        echo -ne "${CYAN}  Waiting for backend... ${elapsed}s${NC}\r"
    done
    echo -e "\n${RED}✗ Backend failed to start within ${BACKEND_HEALTH_TIMEOUT}s${NC}"
    echo -e "${RED}  Check logs: tail -f $BACKEND_LOG${NC}"
    return 1
}

# Create tmux session with first window (backend)
echo -e "${CYAN}Creating tmux session: $SESSION_NAME${NC}"
tmux new-session -d -s "$SESSION_NAME" -n backend

# Start backend in first window
tmux send-keys -t "$SESSION_NAME:backend" "cd $REPO_ROOT/backend && go run . 2>&1 | tee $BACKEND_LOG" C-m

# Wait for backend to be healthy
if ! wait_for_backend; then
    echo -e "${RED}Stopping due to backend failure${NC}"
    tmux kill-session -t "$SESSION_NAME"
    exit 1
fi

# Create second window for frontend
echo -e "\n${YELLOW}[2/2] Starting Expo frontend...${NC}"
tmux new-window -t "$SESSION_NAME" -n frontend
tmux send-keys -t "$SESSION_NAME:frontend" "cd $REPO_ROOT/frontend && npx expo start 2>&1 | tee $FRONTEND_LOG" C-m

# Brief pause for frontend to start
sleep 2
echo -e "${GREEN}✓ Frontend starting...${NC}"

# Print summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}All services started successfully!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${CYAN}tmux session: ${NC}$SESSION_NAME"
echo -e "${CYAN}Attach to session: ${NC}tmux attach -t $SESSION_NAME"
echo ""
echo -e "${CYAN}Windows:${NC}"
echo -e "  0: backend  - Go GraphQL server (http://localhost:8080/graphql)"
echo -e "  1: frontend - Expo dev server"
echo ""
echo -e "${CYAN}Switch windows: ${NC}Ctrl+b then 0/1"
echo -e "${CYAN}Detach session: ${NC}Ctrl+b then d"
echo -e "${CYAN}Stop all services: ${NC}npm run stop"
echo ""
echo -e "${CYAN}Logs:${NC}"
echo -e "  Backend:  tail -f $BACKEND_LOG"
echo -e "  Frontend: tail -f $FRONTEND_LOG"
echo -e "  All logs: npm run logs"
echo ""

# Attach to the session (frontend window)
echo -e "${YELLOW}Attaching to session...${NC}"
echo -e "${CYAN}(Press Ctrl+b then d to detach and keep services running)${NC}"
echo ""
sleep 2
tmux select-window -t "$SESSION_NAME:frontend"
tmux attach -t "$SESSION_NAME"
