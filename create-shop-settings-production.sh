#!/bin/bash

# Script to create default shop settings in production database
# Run this on EC2 server

echo "=========================================="
echo "🔧 Creating Shop Settings in Production"
echo "=========================================="
echo ""

# Set MongoDB URI for production
export MONGODB_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
export MONGODB_DATABASE="cafe_pos"

# Run the create-shop-settings command
cd /home/ubuntu/cafe-pos/backend
go run cmd/create-shop-settings/main.go

echo ""
echo "✅ Done!"
echo ""
