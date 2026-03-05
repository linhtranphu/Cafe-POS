#!/bin/bash

# Script để publish Local Print Bridge lên DockerHub
# Usage: ./publish-to-dockerhub.sh [dockerhub_username] [version]
# Example: ./publish-to-dockerhub.sh myusername v1.0.0
# Example: ./publish-to-dockerhub.sh myusername (uses 'latest')

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DOCKERHUB_USERNAME="${1:-yourusername}"
VERSION="${2:-latest}"
IMAGE_NAME="local-print-bridge"
FULL_IMAGE_NAME="${DOCKERHUB_USERNAME}/${IMAGE_NAME}"

echo -e "${BLUE}=== Publishing Local Print Bridge to DockerHub ===${NC}"
echo ""
echo "DockerHub Username: ${DOCKERHUB_USERNAME}"
echo "Image Name: ${IMAGE_NAME}"
echo "Full Image: ${FULL_IMAGE_NAME}"
echo "Version: ${VERSION}"
echo ""

# Check if username is default
if [ "${DOCKERHUB_USERNAME}" = "yourusername" ]; then
    echo -e "${RED}Error: Please provide your DockerHub username${NC}"
    echo "Usage: ./publish-to-dockerhub.sh <dockerhub_username> [version]"
    echo "Example: ./publish-to-dockerhub.sh myusername v1.0.0"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Check if logged in to DockerHub
echo "Checking DockerHub login..."
if ! docker info 2>/dev/null | grep -q "Username"; then
    echo -e "${YELLOW}Not logged in to DockerHub. Please login:${NC}"
    docker login
    echo ""
fi

# Confirm before proceeding
echo -e "${YELLOW}Ready to build and push:${NC}"
echo "  Image: ${FULL_IMAGE_NAME}:${VERSION}"
if [ "${VERSION}" != "latest" ]; then
    echo "  Also: ${FULL_IMAGE_NAME}:latest"
fi
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# Build image
echo ""
echo -e "${GREEN}Step 1: Building image with Chromium...${NC}"
echo "This will take a few minutes (downloading Chromium and dependencies)..."
echo ""

docker build -t ${FULL_IMAGE_NAME}:${VERSION} .

if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Tag as latest if version is not 'latest'
if [ "${VERSION}" != "latest" ]; then
    echo -e "${GREEN}Step 2: Tagging as latest...${NC}"
    docker tag ${FULL_IMAGE_NAME}:${VERSION} ${FULL_IMAGE_NAME}:latest
    echo -e "${GREEN}✓ Tagged as latest${NC}"
    echo ""
fi

# Show image info
echo -e "${GREEN}Step 3: Image info${NC}"
docker images ${FULL_IMAGE_NAME}
echo ""

# Verify Chromium in image
echo -e "${GREEN}Step 4: Verifying Chromium in image...${NC}"
TEMP_CONTAINER=$(docker run -d ${FULL_IMAGE_NAME}:${VERSION} sleep 10)
if docker exec $TEMP_CONTAINER which chromium-browser > /dev/null 2>&1; then
    CHROME_VERSION=$(docker exec $TEMP_CONTAINER chromium-browser --version 2>&1 || echo "Unknown")
    echo -e "${GREEN}✓ Chromium found: $CHROME_VERSION${NC}"
else
    echo -e "${RED}✗ Chromium not found in image!${NC}"
    docker rm -f $TEMP_CONTAINER
    exit 1
fi
docker rm -f $TEMP_CONTAINER > /dev/null 2>&1
echo ""

# Push version tag
echo -e "${GREEN}Step 5: Pushing ${VERSION} tag to DockerHub...${NC}"
docker push ${FULL_IMAGE_NAME}:${VERSION}

if [ $? -ne 0 ]; then
    echo -e "${RED}Push failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Pushed ${VERSION} tag${NC}"
echo ""

# Push latest tag if version is not 'latest'
if [ "${VERSION}" != "latest" ]; then
    echo -e "${GREEN}Step 6: Pushing latest tag to DockerHub...${NC}"
    docker push ${FULL_IMAGE_NAME}:latest
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Push latest failed!${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ Pushed latest tag${NC}"
    echo ""
fi

# Success summary
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Successfully published to DockerHub!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Image URL: https://hub.docker.com/r/${DOCKERHUB_USERNAME}/${IMAGE_NAME}"
echo ""
echo "Tags pushed:"
echo "  - ${FULL_IMAGE_NAME}:${VERSION}"
if [ "${VERSION}" != "latest" ]; then
    echo "  - ${FULL_IMAGE_NAME}:latest"
fi
echo ""
echo "Users can now pull with:"
echo "  docker pull ${FULL_IMAGE_NAME}:${VERSION}"
if [ "${VERSION}" != "latest" ]; then
    echo "  docker pull ${FULL_IMAGE_NAME}:latest"
fi
echo ""
echo "Or use in docker-compose.yml:"
echo ""
echo "  services:"
echo "    print-bridge:"
echo "      image: ${FULL_IMAGE_NAME}:${VERSION}"
echo "      # ... rest of config"
echo ""
echo "Next steps for users:"
echo "  1. On Linux machine: docker pull ${FULL_IMAGE_NAME}:${VERSION}"
echo "  2. Update docker-compose.yml to use this image"
echo "  3. docker compose up -d"
echo ""
