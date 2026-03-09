#!/bin/bash

set -e

echo "🔨 Rebuilding and deploying frontend..."
echo ""

# Step 1: Build frontend
echo "📦 Step 1: Building frontend..."
cd frontend
rm -rf dist
npm install
npm run build
cd ..

# Step 2: Build Docker image
echo "🐳 Step 2: Building Docker image..."
docker build -t linhtranphu/cafe-pos-frontend:latest -f frontend/Dockerfile frontend/

# Step 3: Push to Docker Hub
echo "📤 Step 3: Pushing to Docker Hub..."
docker push linhtranphu/cafe-pos-frontend:latest

echo ""
echo "✅ Frontend rebuilt and pushed successfully!"
echo ""
echo "📋 Next steps on server:"
echo "1. SSH to server: ssh user@tacafe.store"
echo "2. Pull new image: docker-compose pull frontend"
echo "3. Restart: docker-compose up -d frontend"
echo "4. Clear browser cache: Ctrl+Shift+R"
echo ""
echo "Or run this one-liner on server:"
echo "docker-compose pull frontend && docker-compose up -d frontend"
