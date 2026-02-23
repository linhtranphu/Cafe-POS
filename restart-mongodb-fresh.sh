#!/bin/bash

echo "=== Restarting MongoDB with fresh data ==="
echo ""

# Stop and remove container
echo "1. Stopping and removing MongoDB container..."
docker stop cafe-pos-mongodb
docker rm cafe-pos-mongodb

# Remove volume (optional - will delete all data!)
echo "2. Removing MongoDB volume..."
docker volume rm cafe-pos_mongodb_data 2>/dev/null || true
docker volume rm cafe-pos_mongodb_config 2>/dev/null || true

# Start MongoDB again
echo "3. Starting MongoDB with docker-compose..."
docker-compose up -d mongodb

# Wait for MongoDB to be ready
echo "4. Waiting for MongoDB to be healthy..."
sleep 10

# Check status
echo "5. Checking MongoDB status..."
docker ps | grep mongo

echo ""
echo "✅ MongoDB restarted with fresh data"
echo "Now restart backend to recreate admin user"
