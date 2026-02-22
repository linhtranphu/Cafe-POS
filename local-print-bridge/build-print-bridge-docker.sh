#!/bin/bash

# Build and Push Local Print Bridge Docker Image
# Usage: ./build-print-bridge-docker.sh [version]

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
DOCKER_USERNAME="linhtranphu"
IMAGE_NAME="local-print-bridge"
VERSION=${1:-"1.0.0"}

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🐳 Building Local Print Bridge Docker Image"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Image: ${DOCKER_USERNAME}/${IMAGE_NAME}"
echo "Version: ${VERSION}"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Step 1: Building Docker image...${NC}"
docker build -t ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION} \
             -t ${DOCKER_USERNAME}/${IMAGE_NAME}:latest \
             .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful${NC}"
echo ""

# Check if logged in to Docker Hub
echo -e "${YELLOW}📋 Step 2: Checking Docker Hub login...${NC}"
if ! docker info | grep -q "Username: ${DOCKER_USERNAME}"; then
    echo -e "${YELLOW}⚠️  Not logged in to Docker Hub${NC}"
    echo "Please login to Docker Hub:"
    docker login
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Docker login failed${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}✅ Logged in to Docker Hub${NC}"
echo ""

# Push images
echo -e "${YELLOW}📋 Step 3: Pushing images to Docker Hub...${NC}"
echo "Pushing version ${VERSION}..."
docker push ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Push failed for version ${VERSION}${NC}"
    exit 1
fi

echo "Pushing latest tag..."
docker push ${DOCKER_USERNAME}/${IMAGE_NAME}:latest

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Push failed for latest tag${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Push successful${NC}"
echo ""

# Display summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Build and Push Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 Images pushed:"
echo "  - ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"
echo "  - ${DOCKER_USERNAME}/${IMAGE_NAME}:latest"
echo ""
echo "🚀 To pull and run:"
echo "  docker pull ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"
echo "  docker run -d -p 3001:3001 --name print-bridge ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"
echo ""
echo "📝 Or use docker-compose:"
echo "  cd local-print-bridge"
echo "  docker-compose up -d"
echo ""
