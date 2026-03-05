#!/bin/bash

# Script to rebuild backend with fonts support and redeploy

set -e

echo "=========================================="
echo "🔧 Rebuilding Backend with Fonts Support"
echo "=========================================="
echo ""

# Check Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running"
    exit 1
fi

echo "📦 Building backend image with fonts..."
echo ""

# Build backend with version 2.0.1 (fonts included)
cd backend
docker build --no-cache -t linhtranphu/cafe-pos-backend:2.0.1 .
docker tag linhtranphu/cafe-pos-backend:2.0.1 linhtranphu/cafe-pos-backend:latest

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Backend image built successfully"
else
    echo ""
    echo "❌ Backend build failed"
    exit 1
fi

cd ..

echo ""
echo "📤 Pushing to Docker Hub..."
echo ""

# Login to Docker Hub
docker login

# Push both tags
docker push linhtranphu/cafe-pos-backend:2.0.1
docker push linhtranphu/cafe-pos-backend:latest

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Images pushed successfully"
else
    echo ""
    echo "❌ Push failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ Build Complete!"
echo "=========================================="
echo ""
echo "📝 Next steps on EC2:"
echo "   1. Pull new image: sudo docker-compose pull backend"
echo "   2. Restart backend: sudo docker-compose up -d backend"
echo "   3. Check logs: sudo docker-compose logs -f backend"
echo ""
echo "Images available:"
echo "   - linhtranphu/cafe-pos-backend:2.0.1"
echo "   - linhtranphu/cafe-pos-backend:latest"
echo ""
