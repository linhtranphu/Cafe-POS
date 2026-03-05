#!/bin/bash

# Build and Push Docker Images to Docker Hub
# This script builds frontend, backend, and print bridge images with latest tag

set -e

DOCKER_USERNAME="linhtranphu"
BACKEND_IMAGE="$DOCKER_USERNAME/cafe-pos-backend"
FRONTEND_IMAGE="$DOCKER_USERNAME/cafe-pos-frontend"
PRINT_BRIDGE_IMAGE="$DOCKER_USERNAME/local-print-bridge"

echo "=========================================="
echo "🐳 Docker Hub Build & Push Script"
echo "=========================================="
echo "Tag: latest"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Docker Hub Login
echo "🔐 Logging in to Docker Hub..."
docker login

if [ $? -ne 0 ]; then
    echo "❌ Docker Hub login failed"
    exit 1
fi

echo "✅ Docker Hub login successful"
echo ""

# Build Backend Image
echo "=========================================="
echo "🔨 Building Backend Image..."
echo "=========================================="
echo "Image: $BACKEND_IMAGE:latest"
echo ""

cd backend
docker build --no-cache -t "$BACKEND_IMAGE:latest" .
cd ..

if [ $? -eq 0 ]; then
    echo "✅ Backend image built successfully"
else
    echo "❌ Backend build failed"
    exit 1
fi

echo ""

# Build Frontend Image
echo "=========================================="
echo "🔨 Building Frontend Image..."
echo "=========================================="
echo "Image: $FRONTEND_IMAGE:latest"
echo ""

cd frontend
docker build --no-cache -t "$FRONTEND_IMAGE:latest" .
cd ..

if [ $? -eq 0 ]; then
    echo "✅ Frontend image built successfully"
else
    echo "❌ Frontend build failed"
    exit 1
fi

echo ""

# Build Print Bridge Image
echo "=========================================="
echo "🔨 Building Print Bridge Image..."
echo "=========================================="
echo "Image: $PRINT_BRIDGE_IMAGE:latest"
echo ""

cd local-print-bridge
docker build --no-cache -t "$PRINT_BRIDGE_IMAGE:latest" .
cd ..

if [ $? -eq 0 ]; then
    echo "✅ Print Bridge image built successfully"
    
    # Verify Chromium in image
    echo "🔍 Verifying Chromium in image..."
    TEMP_CONTAINER=$(docker run -d "$PRINT_BRIDGE_IMAGE:latest" sleep 10)
    if docker exec $TEMP_CONTAINER which chromium-browser > /dev/null 2>&1; then
        CHROME_VERSION=$(docker exec $TEMP_CONTAINER chromium-browser --version 2>&1 || echo "Unknown")
        echo "✅ Chromium verified: $CHROME_VERSION"
    else
        echo "⚠️  Warning: Chromium not found in image"
    fi
    docker rm -f $TEMP_CONTAINER > /dev/null 2>&1
else
    echo "❌ Print Bridge build failed"
    exit 1
fi

echo ""

# Push Backend Image
echo "=========================================="
echo "📤 Pushing Backend Image to Docker Hub..."
echo "=========================================="
echo ""

echo "Pushing $BACKEND_IMAGE:latest..."
docker push "$BACKEND_IMAGE:latest"

if [ $? -eq 0 ]; then
    echo "✅ Backend latest pushed successfully"
else
    echo "❌ Backend latest push failed"
    exit 1
fi

echo ""

# Push Frontend Image
echo "=========================================="
echo "📤 Pushing Frontend Image to Docker Hub..."
echo "=========================================="
echo ""

echo "Pushing $FRONTEND_IMAGE:latest..."
docker push "$FRONTEND_IMAGE:latest"

if [ $? -eq 0 ]; then
    echo "✅ Frontend latest pushed successfully"
else
    echo "❌ Frontend latest push failed"
    exit 1
fi

echo ""

# Push Print Bridge Image
echo "=========================================="
echo "📤 Pushing Print Bridge Image to Docker Hub..."
echo "=========================================="
echo ""

echo "Pushing $PRINT_BRIDGE_IMAGE:latest..."
docker push "$PRINT_BRIDGE_IMAGE:latest"

if [ $? -eq 0 ]; then
    echo "✅ Print Bridge latest pushed successfully"
else
    echo "❌ Print Bridge latest push failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ All Done!"
echo "=========================================="
echo ""
echo "📊 Summary:"
echo "  Backend: $BACKEND_IMAGE:latest"
echo "  Frontend: $FRONTEND_IMAGE:latest"
echo "  Print Bridge: $PRINT_BRIDGE_IMAGE:latest"
echo ""
echo "🚀 Images are now available on Docker Hub"
echo "   Ready for deployment!"
echo ""
echo "📝 To deploy on production:"
echo "   Backend/Frontend: docker-compose pull && docker-compose up -d"
echo "   Print Bridge (on local machine):"
echo "     docker pull $PRINT_BRIDGE_IMAGE:latest"
echo "     docker compose up -d"
echo ""
