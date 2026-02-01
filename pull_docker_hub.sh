#!/bin/bash

# Pull Latest Docker Images from Docker Hub and Start Services
# This script pulls the latest images, removes old ones, and starts services on EC2

set -e

DOCKER_USERNAME="linhtranphu"
BACKEND_IMAGE="$DOCKER_USERNAME/cafe-pos-backend:latest"
FRONTEND_IMAGE="$DOCKER_USERNAME/cafe-pos-frontend:latest"
COMPOSE_FILE="docker-compose.yml"

echo "=========================================="
echo "🐳 Docker Hub Pull & Start Script (EC2)"
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

# Remove old images
echo "=========================================="
echo "🗑️  Removing old Docker images..."
echo "=========================================="
echo ""

echo "Removing backend image..."
docker rmi "$BACKEND_IMAGE" 2>/dev/null || echo "ℹ️  Backend image not found locally"

echo "Removing frontend image..."
docker rmi "$FRONTEND_IMAGE" 2>/dev/null || echo "ℹ️  Frontend image not found locally"

echo "✅ Old images removed"
echo ""

# Prune unused images and volumes
echo "=========================================="
echo "🧹 Cleaning up unused Docker resources..."
echo "=========================================="
echo ""

echo "Pruning unused images..."
docker image prune -f --filter "dangling=true" > /dev/null 2>&1 || true

echo "Pruning unused volumes..."
docker volume prune -f > /dev/null 2>&1 || true

echo "✅ Cleanup complete"
echo ""

# Pull latest images
echo "=========================================="
echo "📥 Pulling latest images from Docker Hub..."
echo "=========================================="
echo ""

echo "Pulling backend image: $BACKEND_IMAGE"
docker pull "$BACKEND_IMAGE"

if [ $? -eq 0 ]; then
    echo "✅ Backend image pulled successfully"
else
    echo "❌ Failed to pull backend image"
    exit 1
fi

echo ""

echo "Pulling frontend image: $FRONTEND_IMAGE"
docker pull "$FRONTEND_IMAGE"

if [ $? -eq 0 ]; then
    echo "✅ Frontend image pulled successfully"
else
    echo "❌ Failed to pull frontend image"
    exit 1
fi

echo ""

# Start services
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
echo "  MongoDB:   localhost:27017 (internal only)"
echo ""

echo "📋 Useful Commands:"
echo "  View logs:        docker-compose logs -f"
echo "  View backend logs: docker-compose logs -f backend"
echo "  View frontend logs: docker-compose logs -f frontend"
echo "  View MongoDB logs: docker-compose logs -f mongodb"
echo "  Stop services:    docker-compose down"
echo "  Restart services: docker-compose restart"
echo ""

echo "🚀 Deployment complete!"
echo ""
