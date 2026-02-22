#!/bin/bash

# Stop Local Development Environment
# Stops backend, frontend, Print Bridge, and MongoDB

set -e

echo "=========================================="
echo "🛑 Stopping Local Development Environment"
echo "=========================================="
echo ""

# Stop backend
echo "Stopping backend..."
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    kill -9 $(lsof -t -i:3000) 2>/dev/null || true
    echo "✅ Backend stopped"
else
    echo "ℹ️  Backend not running"
fi

# Stop frontend
echo "Stopping frontend..."
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
    kill -9 $(lsof -t -i:5173) 2>/dev/null || true
    echo "✅ Frontend stopped"
else
    echo "ℹ️  Frontend not running"
fi

# Stop Print Bridge
echo "Stopping Print Bridge..."
if docker ps | grep -q "local-print-bridge"; then
    docker stop local-print-bridge
    echo "✅ Print Bridge stopped"
else
    echo "ℹ️  Print Bridge not running"
fi

# Stop MongoDB
echo "Stopping MongoDB..."
if docker ps | grep -q "cafe-pos-mongodb"; then
    docker-compose -f docker-compose.replica-set.yml down
    echo "✅ MongoDB stopped"
else
    echo "ℹ️  MongoDB not running"
fi

echo ""
echo "=========================================="
echo "✅ All services stopped"
echo "=========================================="
echo ""
echo "To restart: ./restart_local.sh"
echo ""
