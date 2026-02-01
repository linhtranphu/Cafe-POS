#!/bin/bash

# Restart Local Development Environment
# Restarts backend, frontend, and ensures MongoDB is running in Docker

set -e

echo "=========================================="
echo "🔄 Local Development Restart Script"
echo "=========================================="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Check if MongoDB container exists and is running
echo "=========================================="
echo "🗄️  Checking MongoDB..."
echo "=========================================="
echo ""

MONGO_CONTAINER="cafe-pos-mongodb"

if docker ps | grep -q "$MONGO_CONTAINER"; then
    echo "✅ MongoDB is already running"
else
    echo "⚠️  MongoDB is not running. Starting MongoDB..."
    
    # Check if container exists but is stopped
    if docker ps -a | grep -q "$MONGO_CONTAINER"; then
        echo "Starting existing MongoDB container..."
        docker start "$MONGO_CONTAINER"
    else
        echo "Creating and starting new MongoDB container..."
        docker-compose -f docker-compose.local.yml up -d mongodb
    fi
    
    # Wait for MongoDB to be ready
    echo "Waiting for MongoDB to be ready..."
    sleep 3
    
    if docker ps | grep -q "$MONGO_CONTAINER"; then
        echo "✅ MongoDB started successfully"
    else
        echo "❌ Failed to start MongoDB"
        exit 1
    fi
fi

echo ""

# Kill existing backend and frontend processes
echo "=========================================="
echo "🛑 Stopping existing processes..."
echo "=========================================="
echo ""

# Kill backend
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "Stopping backend on port 3000..."
    kill -9 $(lsof -t -i:3000) 2>/dev/null || true
    sleep 1
    echo "✅ Backend stopped"
else
    echo "ℹ️  Backend not running"
fi

# Kill frontend
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "Stopping frontend on port 5173..."
    kill -9 $(lsof -t -i:5173) 2>/dev/null || true
    sleep 1
    echo "✅ Frontend stopped"
else
    echo "ℹ️  Frontend not running"
fi

echo ""

# Start backend
echo "=========================================="
echo "🚀 Starting Backend..."
echo "=========================================="
echo ""

cd backend

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo "❌ go.mod not found in backend directory"
    exit 1
fi

# Set MongoDB URI for local development
export MONGODB_URI="mongodb://admin:password123@localhost:27017"
export MONGODB_DATABASE="cafe_pos"
export JWT_SECRET="your-jwt-secret-key-min-32-chars-long"

# Run backend in background
go run main.go > ../backend.log 2>&1 &
BACKEND_PID=$!

echo "Backend PID: $BACKEND_PID"
echo "Waiting for backend to start..."
sleep 3

# Check if backend is running
if kill -0 $BACKEND_PID 2>/dev/null; then
    echo "✅ Backend started successfully"
else
    echo "❌ Backend failed to start"
    echo "Check backend.log for details"
    cat ../backend.log
    exit 1
fi

cd ..

echo ""

# Start frontend
echo "=========================================="
echo "🚀 Starting Frontend..."
echo "=========================================="
echo ""

cd frontend

# Check if package.json exists
if [ ! -f "package.json" ]; then
    echo "❌ package.json not found in frontend directory"
    exit 1
fi

# Run frontend in background
npm run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!

echo "Frontend PID: $FRONTEND_PID"
echo "Waiting for frontend to start..."
sleep 5

# Check if frontend is running
if kill -0 $FRONTEND_PID 2>/dev/null; then
    echo "✅ Frontend started successfully"
else
    echo "❌ Frontend failed to start"
    echo "Check frontend.log for details"
    cat ../frontend.log
    exit 1
fi

cd ..

echo ""

# Display summary
echo "=========================================="
echo "✅ All Services Started!"
echo "=========================================="
echo ""
echo "📊 Service Status:"
echo "  MongoDB:  ✅ Running on localhost:27017"
echo "  Backend:  ✅ Running on localhost:3000 (PID: $BACKEND_PID)"
echo "  Frontend: ✅ Running on localhost:5173 (PID: $FRONTEND_PID)"
echo ""
echo "🌐 Access Information:"
echo "  Frontend:  http://localhost:5173"
echo "  Backend:   http://localhost:3000"
echo "  MongoDB:   localhost:27017"
echo ""
echo "📋 Logs:"
echo "  Backend:  tail -f backend.log"
echo "  Frontend: tail -f frontend.log"
echo ""
echo "🛑 To stop services:"
echo "  kill $BACKEND_PID  # Stop backend"
echo "  kill $FRONTEND_PID # Stop frontend"
echo "  docker stop cafe-pos-mongodb  # Stop MongoDB"
echo ""
