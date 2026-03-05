#!/bin/bash

# Migration script for Cafe POS v2.0 - PRODUCTION
# Run this script on EC2 server after deploying new Docker images

set -e

echo "=========================================="
echo "🚀 Cafe POS v2.0 - Production Migration"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  Please run with sudo"
    exit 1
fi

# Set production environment variables
export MONGODB_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
export MONGODB_DATABASE="cafe_pos"

echo "📍 Environment: PRODUCTION"
echo "🔗 MongoDB: mongodb://admin:***@localhost:27017/cafe_pos"
echo ""

# Confirm before proceeding
read -p "⚠️  This will modify the PRODUCTION database. Are you sure? (yes/NO) " -r
echo ""
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Migration cancelled"
    exit 1
fi

echo ""
echo "🔄 Starting migration..."
echo ""

# Check if backend container is running
if ! docker ps | grep -q backend; then
    echo "❌ Backend container is not running"
    echo "   Please start services first: docker-compose up -d"
    exit 1
fi

# Get backend container name
BACKEND_CONTAINER=$(docker ps --filter "name=backend" --format "{{.Names}}" | head -1)
echo "🐳 Backend container: $BACKEND_CONTAINER"
echo ""

# Run migration inside Docker container
echo "📦 Running migration in container..."
docker exec -e MONGODB_URI="$MONGODB_URI" -e MONGODB_DATABASE="$MONGODB_DATABASE" \
    $BACKEND_CONTAINER /bin/sh -c "cd /app && go run cmd/migrate-v2.0/main.go"

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "✅ Migration completed successfully!"
    echo "=========================================="
    echo ""
    echo "📝 Next steps:"
    echo "   1. Restart services: docker-compose restart"
    echo "   2. Check logs: docker-compose logs -f backend"
    echo "   3. Test frontend: https://tacafe.store"
    echo ""
else
    echo ""
    echo "=========================================="
    echo "❌ Migration failed!"
    echo "=========================================="
    echo ""
    echo "Please check the error messages above."
    echo "You may need to rollback to previous version."
    exit 1
fi
