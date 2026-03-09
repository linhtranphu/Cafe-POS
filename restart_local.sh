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
        --password 108trannhatduat \
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
            --password 108trannhatduat \
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
        --password 108trannhatduat \
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

# Check and start Print Bridge
echo "=========================================="
echo "🖨️  Checking Print Bridge..."
echo "=========================================="
echo ""

PRINT_BRIDGE_CONTAINER="local-print-bridge"

if docker ps | grep -q "$PRINT_BRIDGE_CONTAINER"; then
    echo "✅ Print Bridge is already running"
else
    echo "⚠️  Print Bridge is not running. Starting..."
    
    # Check if container exists but stopped
    if docker ps -a | grep -q "$PRINT_BRIDGE_CONTAINER"; then
        echo "Starting existing container..."
        docker start "$PRINT_BRIDGE_CONTAINER"
    else
        echo "Creating new container..."
        # Use port mapping instead of host network for macOS compatibility
        docker run -d \
            --name "$PRINT_BRIDGE_CONTAINER" \
            --restart unless-stopped \
            -p 3001:3001 \
            --env-file local-print-bridge/.env \
            linhtranphu/local-print-bridge:latest
    fi
    
    sleep 2
    
    if docker ps | grep -q "$PRINT_BRIDGE_CONTAINER"; then
        echo "✅ Print Bridge started successfully"
    else
        echo "❌ Failed to start Print Bridge"
        echo "Check logs: docker logs $PRINT_BRIDGE_CONTAINER"
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
export MONGODB_URI="mongodb://admin:108trannhatduat@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
export MONGODB_DATABASE="cafe_pos"
export JWT_SECRET="your-jwt-secret-key-min-32-chars-long"

# Build backend
echo "Building backend..."
if go build -v . > ../backend-build.log 2>&1; then
    echo "✅ Backend built successfully"
else
    echo "❌ Backend build failed"
    cat ../backend-build.log
    exit 1
fi

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

# Get local IP address for LAN access
LOCAL_IP=$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || echo "localhost")

# Run frontend in background with --host flag to allow LAN access
npm run dev -- --host > ../frontend.log 2>&1 &
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
echo "  MongoDB:      ✅ Running on localhost:27017 (Replica Set: rs0)"
echo "  Backend:      ✅ Running on localhost:3000 (PID: $BACKEND_PID)"
echo "  Frontend:     ✅ Running on localhost:5173 (PID: $FRONTEND_PID)"
echo "  Print Bridge: ✅ Running on localhost:3001 (Docker)"
echo ""
echo "🌐 Access Information:"
echo "  Frontend (Local):  http://localhost:5173"
echo "  Frontend (LAN):    http://$LOCAL_IP:5173"
echo "  Backend:           http://localhost:3000"
echo "  Print Bridge:      http://localhost:3001/health"
echo "  MongoDB:           mongodb://admin:108trannhatduat@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
echo ""
echo "📱 Access from other devices in LAN:"
echo "  Open browser and go to: http://$LOCAL_IP:5173"
echo ""
echo "📋 Logs:"
echo "  Backend:       tail -f backend.log"
echo "  Frontend:      tail -f frontend.log"
echo "  MongoDB:       docker logs cafe-pos-mongodb"
echo "  Print Bridge:  docker logs local-print-bridge"
echo ""
echo "🛑 To stop services:"
echo "  kill $BACKEND_PID  # Stop backend"
echo "  kill $FRONTEND_PID # Stop frontend"
echo "  docker stop local-print-bridge  # Stop Print Bridge"
echo "  docker-compose -f docker-compose.replica-set.yml down  # Stop MongoDB"
echo ""
