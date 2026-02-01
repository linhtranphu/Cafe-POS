#!/bin/bash

# Build and Push Docker Images to Docker Hub
# This script builds frontend and backend images with no cache and pushes to Docker Hub

set -e

DOCKER_USERNAME="linhtranphu"
BACKEND_IMAGE="$DOCKER_USERNAME/cafe-pos-backend:latest"
FRONTEND_IMAGE="$DOCKER_USERNAME/cafe-pos-frontend:latest"

echo "=========================================="
echo "🐳 Docker Hub Build & Push Script"
echo "=========================================="
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
echo "Image: $BACKEND_IMAGE"
echo ""

cd backend
docker build --no-cache -t "$BACKEND_IMAGE" .
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
echo "Image: $FRONTEND_IMAGE"
echo ""

cd frontend
docker build --no-cache -t "$FRONTEND_IMAGE" .
cd ..

if [ $? -eq 0 ]; then
    echo "✅ Frontend image built successfully"
else
    echo "❌ Frontend build failed"
    exit 1
fi

echo ""

# Push Backend Image
echo "=========================================="
echo "📤 Pushing Backend Image to Docker Hub..."
echo "=========================================="
echo "Image: $BACKEND_IMAGE"
echo ""

docker push "$BACKEND_IMAGE"

if [ $? -eq 0 ]; then
    echo "✅ Backend image pushed successfully"
else
    echo "❌ Backend push failed"
    exit 1
fi

echo ""

# Push Frontend Image
echo "=========================================="
echo "📤 Pushing Frontend Image to Docker Hub..."
echo "=========================================="
echo "Image: $FRONTEND_IMAGE"
echo ""

docker push "$FRONTEND_IMAGE"

if [ $? -eq 0 ]; then
    echo "✅ Frontend image pushed successfully"
else
    echo "❌ Frontend push failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ All Done!"
echo "=========================================="
echo ""
echo "📊 Summary:"
echo "  Backend:  $BACKEND_IMAGE"
echo "  Frontend: $FRONTEND_IMAGE"
echo ""
echo "🚀 Images are now available on Docker Hub"
echo "   Ready for deployment!"
echo ""
