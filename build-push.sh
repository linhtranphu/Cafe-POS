#!/bin/bash

# Build and Push to Docker Hub
# Usage: ./build-push.sh [version]

set -e

VERSION=${1:-latest}
DOCKER_USER="linhtranphu"

echo "🐳 Building and pushing Café POS images to Docker Hub"
echo "======================================================"
echo "Version: $VERSION"
echo "Docker Hub User: $DOCKER_USER"
echo ""

# Login to Docker Hub
echo "📝 Logging in to Docker Hub..."
docker login

# Build Backend
echo ""
echo "🔨 Building backend image..."
cd backend
docker build -t $DOCKER_USER/cafe-pos-backend:$VERSION .
docker tag $DOCKER_USER/cafe-pos-backend:$VERSION $DOCKER_USER/cafe-pos-backend:latest
cd ..

# Build Frontend
echo ""
echo "🔨 Building frontend image..."
cd frontend
docker build -t $DOCKER_USER/cafe-pos-frontend:$VERSION .
docker tag $DOCKER_USER/cafe-pos-frontend:$VERSION $DOCKER_USER/cafe-pos-frontend:latest
cd ..

# Push Backend
echo ""
echo "📤 Pushing backend image..."
docker push $DOCKER_USER/cafe-pos-backend:$VERSION
docker push $DOCKER_USER/cafe-pos-backend:latest

# Push Frontend
echo ""
echo "📤 Pushing frontend image..."
docker push $DOCKER_USER/cafe-pos-frontend:$VERSION
docker push $DOCKER_USER/cafe-pos-frontend:latest

echo ""
echo "✅ Successfully pushed images to Docker Hub!"
echo ""
echo "Images:"
echo "  - $DOCKER_USER/cafe-pos-backend:$VERSION"
echo "  - $DOCKER_USER/cafe-pos-backend:latest"
echo "  - $DOCKER_USER/cafe-pos-frontend:$VERSION"
echo "  - $DOCKER_USER/cafe-pos-frontend:latest"
echo ""
echo "To deploy, run:"
echo "  curl -fsSL https://raw.githubusercontent.com/YOUR-REPO/main/deploy-from-hub.sh | bash"
