#!/bin/bash

set -e

echo "🔧 Fixing 404 frontend error..."
echo ""
echo "This will:"
echo "1. Rebuild frontend with new hash"
echo "2. Build and push Docker image"
echo "3. Update on server"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# Build frontend
echo "📦 Building frontend..."
cd frontend
rm -rf dist
npm run build
cd ..

# Build Docker image with timestamp tag
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
echo "🐳 Building Docker image with tag: $TIMESTAMP"
docker build -t linhtranphu/cafe-pos-frontend:$TIMESTAMP -f frontend/Dockerfile frontend/
docker tag linhtranphu/cafe-pos-frontend:$TIMESTAMP linhtranphu/cafe-pos-frontend:latest

# Push to Docker Hub
echo "📤 Pushing to Docker Hub..."
docker push linhtranphu/cafe-pos-frontend:$TIMESTAMP
docker push linhtranphu/cafe-pos-frontend:latest

echo ""
echo "✅ Done! Now run on server:"
echo ""
echo "ssh user@tacafe.store"
echo "cd /path/to/project"
echo "docker-compose pull frontend"
echo "docker-compose up -d frontend"
echo ""
echo "Then clear browser cache: Ctrl+Shift+R or Cmd+Shift+R"
