#!/bin/bash

# Migration script for Cafe POS v2.0
# This script updates the database structure without losing existing data

set -e

echo "=========================================="
echo "🚀 Cafe POS - Migration to v2.0"
echo "=========================================="
echo ""

# Check if running on production or local
if [ -f ".env.production" ]; then
    echo "📍 Environment: PRODUCTION"
    export $(cat .env.production | grep -v '^#' | xargs)
elif [ -f ".env" ]; then
    echo "📍 Environment: LOCAL"
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  No .env file found, using defaults"
fi

# Display MongoDB connection info (hide password)
MONGO_URI_DISPLAY=$(echo $MONGODB_URI | sed 's/:.*@/:***@/')
echo "🔗 MongoDB: $MONGO_URI_DISPLAY"
echo ""

# Confirm before proceeding
read -p "⚠️  This will modify the database structure. Continue? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Migration cancelled"
    exit 1
fi

echo ""
echo "🔄 Starting migration..."
echo ""

# Run migration
cd backend
go run cmd/migrate-v2.0/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "✅ Migration completed successfully!"
    echo "=========================================="
    echo ""
    echo "📝 Next steps:"
    echo "   1. Restart backend: docker-compose restart backend"
    echo "   2. Clear browser cache and reload frontend"
    echo "   3. Test shop settings: http://localhost:5173/#/print-management"
    echo ""
else
    echo ""
    echo "=========================================="
    echo "❌ Migration failed!"
    echo "=========================================="
    echo ""
    echo "Please check the error messages above."
    exit 1
fi
