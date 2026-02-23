#!/bin/bash

# Quick script to restart both backend and frontend servers
# Usage: ./RESTART_SERVERS.sh

echo "🔄 Restarting Cafe POS Servers"
echo "================================"
echo ""

# Function to kill process on port
kill_port() {
    local port=$1
    local pid=$(lsof -ti:$port)
    if [ ! -z "$pid" ]; then
        echo "Killing process on port $port (PID: $pid)"
        kill -9 $pid 2>/dev/null
        sleep 1
    fi
}

# Kill existing processes
echo "🛑 Stopping existing servers..."
kill_port 8080  # Backend
kill_port 5173  # Frontend
echo ""

# Start backend
echo "🚀 Starting backend on port 8080..."
cd backend
go run main.go &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"
cd ..
sleep 3
echo ""

# Start frontend
echo "🚀 Starting frontend on port 5173..."
cd frontend
npm run dev &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"
cd ..
sleep 2
echo ""

echo "✅ Servers started!"
echo ""
echo "Backend:  http://localhost:8080"
echo "Frontend: http://localhost:5173"
echo ""
echo "To stop servers:"
echo "  kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "Or use Ctrl+C to stop this script"
echo ""

# Wait for user interrupt
wait
