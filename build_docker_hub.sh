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

# ── Ask which images to build/push ────────────────────────────────────────────
echo "Select images to build & push:"
echo "  [1] Backend only"
echo "  [2] Frontend only"
echo "  [3] Print Bridge only"
echo "  [4] Backend + Frontend"
echo "  [5] All (Backend + Frontend + Print Bridge)"
echo ""
read -r -p "Choice [1-5]: " CHOICE

BUILD_BACKEND=false
BUILD_FRONTEND=false
BUILD_PRINT_BRIDGE=false

case "$CHOICE" in
  1) BUILD_BACKEND=true ;;
  2) BUILD_FRONTEND=true ;;
  3) BUILD_PRINT_BRIDGE=true ;;
  4) BUILD_BACKEND=true; BUILD_FRONTEND=true ;;
  5) BUILD_BACKEND=true; BUILD_FRONTEND=true; BUILD_PRINT_BRIDGE=true ;;
  *)
    echo "❌ Invalid choice. Exiting."
    exit 1
    ;;
esac

echo ""
echo "Will build/push:"
[[ $BUILD_BACKEND == true ]]      && echo "  ✔ Backend"
[[ $BUILD_FRONTEND == true ]]     && echo "  ✔ Frontend"
[[ $BUILD_PRINT_BRIDGE == true ]] && echo "  ✔ Print Bridge"
echo ""
read -r -p "Proceed? [y/N] " CONFIRM
if [[ "$(echo "$CONFIRM" | tr '[:upper:]' '[:lower:]')" != "y" ]]; then
    echo "Aborted."
    exit 0
fi
echo ""

# Docker Hub Login
echo "🔐 Logging in to Docker Hub..."
docker login

if [ $? -ne 0 ]; then
    echo "❌ Docker Hub login failed"
    exit 1
fi

echo "✅ Docker Hub login successful"
echo ""

# ── Build Backend ──────────────────────────────────────────────────────────────
if [[ $BUILD_BACKEND == true ]]; then
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
fi

# ── Build Frontend ─────────────────────────────────────────────────────────────
if [[ $BUILD_FRONTEND == true ]]; then
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
fi

# ── Build Print Bridge ─────────────────────────────────────────────────────────
if [[ $BUILD_PRINT_BRIDGE == true ]]; then
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
fi

# ── Push Backend ───────────────────────────────────────────────────────────────
if [[ $BUILD_BACKEND == true ]]; then
    echo "=========================================="
    echo "📤 Pushing Backend Image to Docker Hub..."
    echo "=========================================="

    docker push "$BACKEND_IMAGE:latest"

    if [ $? -eq 0 ]; then
        echo "✅ Backend pushed successfully"
    else
        echo "❌ Backend push failed"
        exit 1
    fi
    echo ""
fi

# ── Push Frontend ──────────────────────────────────────────────────────────────
if [[ $BUILD_FRONTEND == true ]]; then
    echo "=========================================="
    echo "📤 Pushing Frontend Image to Docker Hub..."
    echo "=========================================="

    docker push "$FRONTEND_IMAGE:latest"

    if [ $? -eq 0 ]; then
        echo "✅ Frontend pushed successfully"
    else
        echo "❌ Frontend push failed"
        exit 1
    fi
    echo ""
fi

# ── Push Print Bridge ──────────────────────────────────────────────────────────
if [[ $BUILD_PRINT_BRIDGE == true ]]; then
    echo "=========================================="
    echo "📤 Pushing Print Bridge Image to Docker Hub..."
    echo "=========================================="

    docker push "$PRINT_BRIDGE_IMAGE:latest"

    if [ $? -eq 0 ]; then
        echo "✅ Print Bridge pushed successfully"
    else
        echo "❌ Print Bridge push failed"
        exit 1
    fi
    echo ""
fi

echo ""

# Cleanup Docker to free up space
echo "=========================================="
echo "🧹 Cleaning up Docker..."
echo "=========================================="
echo ""

echo "Removing dangling images..."
docker image prune -f

echo ""
echo "Removing unused build cache..."
docker builder prune -f

echo ""
echo "Removing stopped containers..."
docker container prune -f

echo ""
echo "📊 Docker disk usage after cleanup:"
docker system df

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
echo "💡 Tip: If you need more space, run:"
echo "   docker system prune -a --volumes"
echo "   (This will remove ALL unused images and volumes)"
echo ""
echo "📝 To deploy on production:"
echo "   Backend/Frontend: docker-compose pull && docker-compose up -d"
echo "   Print Bridge (on local machine):"
echo "     docker pull $PRINT_BRIDGE_IMAGE:latest"
echo "     docker compose up -d"
echo ""
