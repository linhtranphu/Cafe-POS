#!/bin/bash

# Script to check if cancel order endpoint is available for waiter

echo "=========================================="
echo "🔍 Checking Cancel Order Endpoint"
echo "=========================================="
echo ""

# Check if backend is running
echo "1. Checking if backend is running..."
if docker ps | grep -q cafe-pos-backend; then
    echo "✅ Backend container is running"
else
    echo "❌ Backend container is not running"
    exit 1
fi

echo ""
echo "2. Checking backend logs for waiter routes..."
docker logs cafe-pos-backend 2>&1 | grep -i "waiter" | tail -5

echo ""
echo "3. Checking if cancel endpoint is registered..."
docker logs cafe-pos-backend 2>&1 | grep -i "cancel" | tail -5

echo ""
echo "=========================================="
echo "📝 Recommendations"
echo "=========================================="
echo ""
echo "If you don't see the cancel endpoint registered:"
echo "1. Rebuild backend:"
echo "   cd backend"
echo "   docker-compose build --no-cache"
echo "   docker-compose up -d"
echo ""
echo "2. Or restart the entire stack:"
echo "   docker-compose down"
echo "   docker-compose up -d --build"
echo ""
echo "3. Check backend version:"
echo "   docker exec cafe-pos-backend cat /app/main.go | grep 'waiter.POST.*cancel'"
echo ""
