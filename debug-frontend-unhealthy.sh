#!/bin/bash

# Debug Frontend Unhealthy Container
# This script helps diagnose why frontend container is unhealthy

set -e

echo "🔍 Debugging Frontend Container Health"
echo "======================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CONTAINER_NAME="cafe-pos-frontend"

# Check if container exists
echo -e "${BLUE}1️⃣ Checking if container exists...${NC}"
if ! docker ps -a | grep -q "$CONTAINER_NAME"; then
  echo -e "${RED}❌ Container $CONTAINER_NAME not found${NC}"
  exit 1
fi
echo -e "${GREEN}✅ Container found${NC}"
echo ""

# Check container status
echo -e "${BLUE}2️⃣ Container Status:${NC}"
docker ps -a | grep "$CONTAINER_NAME"
echo ""

# Check container logs (last 50 lines)
echo -e "${BLUE}3️⃣ Container Logs (last 50 lines):${NC}"
echo "-----------------------------------"
docker logs --tail 50 "$CONTAINER_NAME"
echo ""

# Check healthcheck logs
echo -e "${BLUE}4️⃣ Healthcheck Status:${NC}"
docker inspect --format='{{json .State.Health}}' "$CONTAINER_NAME" | jq '.'
echo ""

# Test healthcheck command manually
echo -e "${BLUE}5️⃣ Testing Healthcheck Command Manually:${NC}"
echo "Command: wget --quiet --tries=1 --spider http://localhost"
docker exec "$CONTAINER_NAME" wget --quiet --tries=1 --spider http://localhost && \
  echo -e "${GREEN}✅ Healthcheck command succeeded${NC}" || \
  echo -e "${RED}❌ Healthcheck command failed${NC}"
echo ""

# Check if nginx is running
echo -e "${BLUE}6️⃣ Checking if nginx is running:${NC}"
docker exec "$CONTAINER_NAME" ps aux | grep nginx || echo -e "${RED}❌ nginx not found${NC}"
echo ""

# Check nginx configuration
echo -e "${BLUE}7️⃣ Testing nginx configuration:${NC}"
docker exec "$CONTAINER_NAME" nginx -t 2>&1 || echo -e "${RED}❌ nginx config test failed${NC}"
echo ""

# Check port 80 is listening
echo -e "${BLUE}8️⃣ Checking if port 80 is listening:${NC}"
docker exec "$CONTAINER_NAME" netstat -tuln | grep :80 || \
  echo -e "${YELLOW}⚠️ Port 80 not listening (netstat not available, trying ss)${NC}"
docker exec "$CONTAINER_NAME" ss -tuln | grep :80 || \
  echo -e "${RED}❌ Port 80 not listening${NC}"
echo ""

# Check files in /usr/share/nginx/html
echo -e "${BLUE}9️⃣ Checking files in nginx html directory:${NC}"
docker exec "$CONTAINER_NAME" ls -lah /usr/share/nginx/html/ || \
  echo -e "${RED}❌ Cannot list files${NC}"
echo ""

# Test HTTP request from inside container
echo -e "${BLUE}🔟 Testing HTTP request from inside container:${NC}"
docker exec "$CONTAINER_NAME" wget -O- http://localhost 2>&1 | head -20 || \
  echo -e "${RED}❌ HTTP request failed${NC}"
echo ""

# Test HTTP request from host
echo -e "${BLUE}1️⃣1️⃣ Testing HTTP request from host:${NC}"
curl -I http://localhost 2>&1 | head -10 || \
  echo -e "${RED}❌ Cannot reach from host${NC}"
echo ""

# Check container resource usage
echo -e "${BLUE}1️⃣2️⃣ Container Resource Usage:${NC}"
docker stats --no-stream "$CONTAINER_NAME"
echo ""

# Summary
echo -e "${YELLOW}📋 Summary:${NC}"
echo "-----------------------------------"
echo "Container: $CONTAINER_NAME"
echo "Status: $(docker inspect --format='{{.State.Status}}' $CONTAINER_NAME)"
echo "Health: $(docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME)"
echo ""

echo -e "${YELLOW}💡 Common Issues:${NC}"
echo "1. nginx not started - Check logs for startup errors"
echo "2. Port 80 not listening - nginx config issue"
echo "3. Files not built - Check if dist files exist"
echo "4. wget not installed - Healthcheck command fails"
echo "5. nginx config error - Run 'nginx -t' to test"
echo ""

echo -e "${YELLOW}🔧 Quick Fixes:${NC}"
echo "1. Restart container: docker restart $CONTAINER_NAME"
echo "2. Check logs: docker logs -f $CONTAINER_NAME"
echo "3. Exec into container: docker exec -it $CONTAINER_NAME sh"
echo "4. Rebuild image: docker-compose -f docker-compose.hub.yml build frontend"
echo "5. Remove and recreate: docker-compose -f docker-compose.hub.yml up -d --force-recreate frontend"
echo ""

echo "✅ Debug complete!"
