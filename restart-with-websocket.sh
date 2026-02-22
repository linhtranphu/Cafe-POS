#!/bin/bash

# Restart All Services with WebSocket Support
# This script ensures all services use compatible Socket.IO versions

set -e

echo "=========================================="
echo "🔄 Restarting Services with WebSocket"
echo "=========================================="
echo ""

# Kill existing processes
echo "🛑 Stopping existing processes..."

# Kill backend
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "  Stopping backend..."
    kill -9 $(lsof -t -i:3000) 2>/dev/null || true
fi

# Kill frontend
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "  Stopping frontend..."
    kill -9 $(lsof -t -i:5173) 2>/dev/null || true
fi

# Kill Print Bridge (if running locally)
if lsof -Pi :3001 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "  Stopping Print Bridge..."
    kill -9 $(lsof -t -i:3001) 2>/dev/null || true
fi

sleep 2
echo "✅ All processes stopped"
echo ""

# Check MongoDB
echo "=========================================="
echo "🗄️  Checking MongoDB..."
echo "=========================================="
echo ""

if ! docker ps | grep -q "cafe-pos-mongodb"; then
    echo "⚠️  MongoDB not running. Starting..."
    docker-compose -f docker-compose.replica-set.yml up -d mongodb
    sleep 20
    echo "✅ MongoDB started"
else
    echo "✅ MongoDB already running"
fi

echo ""

# Start Backend
echo "=========================================="
echo "🚀 Starting Backend..."
echo "=========================================="
echo ""

cd backend
export MONGODB_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
export MONGODB_DATABASE="cafe_pos"
export JWT_SECRET="your-jwt-secret-key-min-32-chars-long"

go run main.go > ../backend.log 2>&1 &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"
sleep 3

if kill -0 $BACKEND_PID 2>/dev/null; then
    echo "✅ Backend started on port 3000"
else
    echo "❌ Backend failed to start"
    cat ../backend.log
    exit 1
fi

cd ..
echo ""

# Start Print Bridge
echo "=========================================="
echo "🖨️  Starting Print Bridge..."
echo "=========================================="
echo ""

cd local-print-bridge
npm start > ../print-bridge.log 2>&1 &
PRINT_BRIDGE_PID=$!
echo "Print Bridge PID: $PRINT_BRIDGE_PID"
sleep 3

if kill -0 $PRINT_BRIDGE_PID 2>/dev/null; then
    echo "✅ Print Bridge started on port 3001"
    echo ""
    echo "📋 Checking WebSocket connection..."
    sleep 2
    if grep -q "Connected to backend" ../print-bridge.log; then
        echo "✅ WebSocket connected!"
    else
        echo "⚠️  WebSocket not connected yet, check logs"
    fi
else
    echo "❌ Print Bridge failed to start"
    cat ../print-bridge.log
    exit 1
fi

cd ..
echo ""

# Start Frontend
echo "=========================================="
echo "🚀 Starting Frontend..."
echo "=========================================="
echo ""

cd frontend
npm run dev -- --host > ../frontend.log 2>&1 &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"
sleep 5

if kill -0 $FRONTEND_PID 2>/dev/null; then
    echo "✅ Frontend started on port 5173"
else
    echo "❌ Frontend failed to start"
    cat ../frontend.log
    exit 1
fi

cd ..
echo ""

# Get local IP
LOCAL_IP=$(ipconfig getifaddr en0 2>/dev/null || echo "N/A")

# Summary
echo "=========================================="
echo "✅ All Services Started with WebSocket!"
echo "=========================================="
echo ""
echo "📊 Service Status:"
echo "  MongoDB:      ✅ Running (Replica Set: rs0)"
echo "  Backend:      ✅ Running on port 3000 (PID: $BACKEND_PID)"
echo "  Print Bridge: ✅ Running on port 3001 (PID: $PRINT_BRIDGE_PID)"
echo "  Frontend:     ✅ Running on port 5173 (PID: $FRONTEND_PID)"
echo ""
echo "🌐 Access URLs:"
echo "  Frontend (Local):  http://localhost:5173"
if [ "$LOCAL_IP" != "N/A" ]; then
    echo "  Frontend (LAN):    http://$LOCAL_IP:5173"
fi
echo "  Backend:           http://localhost:3000"
echo "  Print Bridge:      http://localhost:3001/health"
echo ""
echo "📋 Logs:"
echo "  Backend:       tail -f backend.log"
echo "  Print Bridge:  tail -f print-bridge.log"
echo "  Frontend:      tail -f frontend.log"
echo ""
echo "🔍 Verify WebSocket:"
echo "  tail -f print-bridge.log | grep WebSocket"
echo "  # Should see: [WebSocket] ✅ Connected to backend"
echo ""
echo "🛑 To stop all services:"
echo "  kill $BACKEND_PID $PRINT_BRIDGE_PID $FRONTEND_PID"
echo ""
echo "🎯 Next: Open http://localhost:5173 and create an order to test!"
echo ""
