#!/bin/bash

# Quick Frontend Redeploy Script
# Rebuilds and redeploys only the frontend to fix cache issues

set -e

DOCKER_USERNAME="linhtranphu"
FRONTEND_IMAGE="$DOCKER_USERNAME/cafe-pos-frontend"

echo "=========================================="
echo "🔄 Frontend Quick Redeploy"
echo "=========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi

# Build Frontend
echo "🔨 Building Frontend..."
cd frontend
docker build --no-cache -t "$FRONTEND_IMAGE:latest" .
cd ..

if [ $? -ne 0 ]; then
    echo "❌ Frontend build failed"
    exit 1
fi

echo "✅ Frontend built successfully"
echo ""

# Push to Docker Hub
echo "📤 Pushing to Docker Hub..."
docker push "$FRONTEND_IMAGE:latest"

if [ $? -ne 0 ]; then
    echo "❌ Push failed"
    exit 1
fi

echo "✅ Pushed successfully"
echo ""

# Cleanup
echo "🧹 Cleaning up..."
docker image prune -f
docker builder prune -f

echo ""
echo "=========================================="
echo "✅ Frontend Redeployed!"
echo "=========================================="
echo ""
echo "📝 Next steps on production server:"
echo "   1. docker-compose pull"
echo "   2. docker-compose up -d frontend"
echo "   3. docker system prune -f"
echo ""
echo "💡 Users should hard refresh browser:"
echo "   - Chrome/Edge: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)"
echo "   - Safari: Cmd+Option+R"
echo ""
