#!/bin/bash

echo "=== Restarting Backend ==="
echo ""

# Kill existing backend process
echo "1. Stopping existing backend..."
pkill -f "go run main.go" || echo "   No existing process found"
sleep 2

# Start backend in background
echo "2. Starting backend..."
cd backend
nohup go run main.go > ../backend.log 2>&1 &
BACKEND_PID=$!

echo "   Backend started with PID: $BACKEND_PID"
echo ""

# Wait a bit for backend to start
echo "3. Waiting for backend to start..."
sleep 3

# Check if backend is running
if ps -p $BACKEND_PID > /dev/null; then
    echo "   ✅ Backend is running"
    echo ""
    echo "Checking logs..."
    tail -20 ../backend.log
else
    echo "   ❌ Backend failed to start"
    echo ""
    echo "Error logs:"
    cat ../backend.log
    exit 1
fi

echo ""
echo "=== Backend Restarted Successfully ==="
echo ""
echo "To view logs: tail -f backend.log"
echo "To stop: pkill -f 'go run main.go'"
