#!/bin/bash

# Local Print Bridge - Docker Start Script

set -e

echo "=========================================="
echo "Local Print Bridge - Docker Deployment"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found${NC}"
    echo "Creating .env from .env.docker template..."
    cp .env.docker .env
    echo -e "${GREEN}✓ .env file created${NC}"
    echo ""
    echo -e "${YELLOW}IMPORTANT: Please edit .env file and update:${NC}"
    echo "  - BACKEND_URL (your EC2 domain)"
    echo "  - DEFAULT_BILL_PRINTER_IP (your printer IP)"
    echo "  - DEFAULT_LABEL_PRINTER_IP (your printer IP)"
    echo ""
    read -p "Press Enter after updating .env file..."
fi

# Load environment variables
source .env

# Validate required variables
if [ "$BACKEND_URL" = "https://your-ec2-domain.com" ]; then
    echo -e "${RED}✗ Please update BACKEND_URL in .env file${NC}"
    exit 1
fi

echo "Configuration:"
echo "  Backend URL: $BACKEND_URL"
echo "  Bill Printer: $DEFAULT_BILL_PRINTER_IP"
echo "  Label Printer: $DEFAULT_LABEL_PRINTER_IP"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is not installed${NC}"
    echo "Please install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}✗ Docker Compose is not installed${NC}"
    echo "Please install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "Building Docker image..."
docker-compose build

echo ""
echo "Starting Local Print Bridge..."
docker-compose up -d

echo ""
echo "Waiting for service to start..."
sleep 3

# Check if container is running
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✓ Local Print Bridge is running${NC}"
    echo ""
    
    # Show logs
    echo "Recent logs:"
    docker-compose logs --tail=20
    
    echo ""
    echo "=========================================="
    echo "Service Information"
    echo "=========================================="
    echo "  Status: Running"
    echo "  URL: http://localhost:3001"
    echo "  Health: http://localhost:3001/health"
    echo ""
    echo "Useful Commands:"
    echo "  View logs:    docker-compose logs -f"
    echo "  Stop service: docker-compose stop"
    echo "  Restart:      docker-compose restart"
    echo "  Remove:       docker-compose down"
    echo ""
    echo "Testing:"
    echo "  curl http://localhost:3001/health"
    echo ""
else
    echo -e "${RED}✗ Failed to start Local Print Bridge${NC}"
    echo ""
    echo "Logs:"
    docker-compose logs
    exit 1
fi
