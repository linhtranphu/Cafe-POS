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
echo "🗄️  Checking MongoDB Replica Set..."
echo "=========================================="
echo ""

MONGO_CONTAINER="cafe-pos-mongodb"

if docker ps | grep -q "$MONGO_CONTAINER"; then
    echo "✅ MongoDB is already running"
    
    # Check if replica set is initialized
    echo "Checking replica set status..."
    RS_STATUS=$(docker exec "$MONGO_CONTAINER" mongosh \
        --username admin \
        --password password123 \
        --authenticationDatabase admin \
        --quiet \
        --eval "try { rs.status().ok } catch(e) { 0 }" 2>/dev/null || echo "0")
    
    if [ "$RS_STATUS" = "1" ]; then
        echo "✅ Replica set is active"
    else
        echo "⚠️  Replica set not initialized. Restarting with replica set configuration..."
        docker-compose down
        docker-compose -f docker-compose.replica-set.yml up -d mongodb
        
        echo "Waiting for MongoDB to start..."
        sleep 20
        
        echo "Initializing replica set..."
        docker exec "$MONGO_CONTAINER" mongosh \
            --username admin \
            --password password123 \
            --authenticationDatabase admin \
            --eval "
            try {
                rs.status();
                print('Replica set already initialized');
            } catch(e) {
                rs.initiate({
                    _id: 'rs0',
                    members: [{ _id: 0, host: 'localhost:27017' }]
                });
                print('Replica set initialized');
            }
            " > /dev/null 2>&1
        
        echo "Waiting for replica set to stabilize..."
        sleep 10
        echo "✅ Replica set ready"
    fi
else
    echo "⚠️  MongoDB is not running. Starting MongoDB with replica set..."
    
    # Stop any existing containers
    docker-compose down 2>/dev/null || true
    
    # Start with replica set configuration
    docker-compose -f docker-compose.replica-set.yml up -d mongodb
    
    # Wait for MongoDB to be ready
    echo "Waiting for MongoDB to start..."
    sleep 20
    
    # Initialize replica set
    echo "Initializing replica set..."
    docker exec "$MONGO_CONTAINER" mongosh \
        --username admin \
        --password password123 \
        --authenticationDatabase admin \
        --eval "
        try {
            rs.status();
            print('Replica set already initialized');
        } catch(e) {
            rs.initiate({
                _id: 'rs0',
                members: [{ _id: 0, host: 'localhost:27017' }]
            });
            print('Replica set initialized');
        }
        " > /dev/null 2>&1
    
    echo "Waiting for replica set to stabilize..."
    sleep 10
    
    if docker ps | grep -q "$MONGO_CONTAINER"; then
        echo "✅ MongoDB replica set started successfully"
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

# Set MongoDB URI for local development with replica set
export MONGODB_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
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
echo "  MongoDB:  ✅ Running on localhost:27017 (Replica Set: rs0)"
echo "  Backend:  ✅ Running on localhost:3000 (PID: $BACKEND_PID)"
echo "  Frontend: ✅ Running on localhost:5173 (PID: $FRONTEND_PID)"
echo ""
echo "🌐 Access Information:"
echo "  Frontend:  http://localhost:5173"
echo "  Backend:   http://localhost:3000"
echo "  MongoDB:   mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
echo ""
echo "📋 Logs:"
echo "  Backend:  tail -f backend.log"
echo "  Frontend: tail -f frontend.log"
echo "  MongoDB:  docker logs cafe-pos-mongodb"
echo ""
echo "🛑 To stop services:"
echo "  kill $BACKEND_PID  # Stop backend"
echo "  kill $FRONTEND_PID # Stop frontend"
echo "  docker-compose -f docker-compose.replica-set.yml down  # Stop MongoDB"
echo ""
