#!/bin/bash

# Script to rebuild backend with cancel order feature for waiter

set -e

echo "=========================================="
echo "🔨 Rebuilding Backend with Cancel Feature"
echo "=========================================="
echo ""

# Confirm before proceeding
read -p "This will rebuild and restart the backend. Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""
echo "1. Stopping backend container..."
docker-compose stop backend

echo ""
echo "2. Removing old backend container..."
docker-compose rm -f backend

echo ""
echo "3. Rebuilding backend (no cache)..."
cd backend
docker build --no-cache -t cafe-pos-backend:latest .
cd ..

echo ""
echo "4. Starting backend..."
docker-compose up -d backend

echo ""
echo "5. Waiting for backend to start..."
sleep 5

echo ""
echo "6. Checking backend health..."
docker ps | grep cafe-pos-backend

echo ""
echo "7. Checking backend logs..."
docker logs cafe-pos-backend 2>&1 | tail -20

echo ""
echo "=========================================="
echo "✅ Backend Rebuild Complete"
echo "=========================================="
echo ""
echo "📝 Next steps:"
echo "1. Test cancel order feature in the app"
echo "2. Check logs if there are issues:"
echo "   docker logs -f cafe-pos-backend"
echo ""
