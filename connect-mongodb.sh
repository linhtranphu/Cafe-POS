#!/bin/bash

# MongoDB Connection Script
# Quick access to MongoDB running in Docker

# Load environment variables from .env if exists
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
elif [ -f ".env.production" ]; then
    export $(cat .env.production | grep -v '^#' | grep -v '^$' | xargs)
fi

# Get MongoDB credentials from environment or use defaults
MONGO_USER="${MONGO_INITDB_ROOT_USERNAME:-admin}"
MONGO_PASS="${MONGO_INITDB_ROOT_PASSWORD:-password}"
MONGO_DB="${MONGODB_DATABASE:-cafe_pos}"

echo "🔗 Connecting to MongoDB..."
echo "   User: $MONGO_USER"
echo "   Database: $MONGO_DB"
echo ""

# Connect to MongoDB
docker exec -it cafe-pos-mongodb mongosh \
    -u "$MONGO_USER" \
    -p "$MONGO_PASS" \
    --authenticationDatabase admin \
    "$MONGO_DB"
