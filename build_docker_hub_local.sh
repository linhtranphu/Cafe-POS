#!/bin/bash

# Build and Run Docker Images Locally
# This script builds frontend and backend images locally and runs them with docker-compose

set -e

DOCKER_USERNAME="linhtranphu"
BACKEND_IMAGE="$DOCKER_USERNAME/cafe-pos-backend:latest"
FRONTEND_IMAGE="$DOCKER_USERNAME/cafe-pos-frontend:latest"
COMPOSE_FILE="docker-compose.local.yml"

echo "=========================================="
echo "🐳 Docker Local Build & Run Script"
echo "=========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Stop running containers
echo "=========================================="
echo "🛑 Stopping running containers..."
echo "=========================================="
echo ""

if docker-compose -f "$COMPOSE_FILE" ps | grep -q "Up"; then
    echo "Stopping services..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "✅ Services stopped"
else
    echo "ℹ️  No running services found"
fi

echo ""

# Clean up Docker resources
echo "=========================================="
echo "🧹 Cleaning up Docker resources..."
echo "=========================================="
echo ""

echo "Removing unused containers..."
docker container prune -f > /dev/null 2>&1

echo "Removing unused images..."
docker image prune -f > /dev/null 2>&1

echo "✅ Docker cleanup complete"
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

# Restart Docker daemon
echo "=========================================="
echo "🔄 Restarting Docker daemon..."
echo "=========================================="
echo ""



# Start services with docker-compose
echo "=========================================="
echo "🚀 Starting services with docker-compose..."
echo "=========================================="
echo ""

docker-compose -f "$COMPOSE_FILE" up -d

if [ $? -eq 0 ]; then
    echo "✅ Services started successfully"
else
    echo "❌ Failed to start services"
    exit 1
fi

echo ""

# Wait for services to be healthy
echo "=========================================="
echo "⏳ Waiting for services to be healthy..."
echo "=========================================="
echo ""

sleep 5

# Check service status
echo "Checking service status..."
docker-compose -f "$COMPOSE_FILE" ps

echo ""

# Display access information
echo "=========================================="
echo "✅ All Done!"
echo "=========================================="
echo ""
echo "📊 Services Status:"
docker-compose -f "$COMPOSE_FILE" ps --services

echo ""
echo "🌐 Access Information:"
echo "  Frontend:  http://localhost"
echo "  Backend:   http://localhost:3000"
echo "  MongoDB:   localhost:27017"
echo ""

echo "📋 Useful Commands:"
echo "  View logs:        docker-compose -f $COMPOSE_FILE logs -f"
echo "  View backend logs: docker-compose -f $COMPOSE_FILE logs -f backend"
echo "  View frontend logs: docker-compose -f $COMPOSE_FILE logs -f frontend"
echo "  View MongoDB logs: docker-compose -f $COMPOSE_FILE logs -f mongodb"
echo "  Stop services:    docker-compose -f $COMPOSE_FILE down"
echo "  Restart services: docker-compose -f $COMPOSE_FILE restart"
echo ""

echo "🚀 Local deployment complete!"
echo ""
