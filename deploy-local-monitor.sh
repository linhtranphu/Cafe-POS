#!/bin/bash

# Script để deploy local và monitor resource usage
# Mục đích: Tìm nguyên nhân server bị chết trên EC2

set -e

echo "=========================================="
echo "🚀 Deploy Local với Resource Monitoring"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}⚠️  Dừng monitoring...${NC}"
    if [ ! -z "$MONITOR_PID" ]; then
        kill $MONITOR_PID 2>/dev/null || true
    fi
}

trap cleanup EXIT

# Step 1: Stop existing containers
echo -e "${YELLOW}📦 Dừng containers cũ...${NC}"
docker-compose -f docker-compose.local.yml down 2>/dev/null || true
echo ""

# Step 2: Build backend
echo -e "${GREEN}🔨 Build Backend...${NC}"
cd backend
docker build -t cafe-pos-backend:local . --no-cache
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Backend build failed!${NC}"
    exit 1
fi
cd ..
echo ""

# Step 3: Build frontend
echo -e "${GREEN}🔨 Build Frontend...${NC}"
cd frontend
docker build -t cafe-pos-frontend:local . --no-cache
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Frontend build failed!${NC}"
    exit 1
fi
cd ..
echo ""

# Step 4: Start services
echo -e "${GREEN}🚀 Starting services...${NC}"
docker-compose -f docker-compose.local.yml up -d
echo ""

# Step 5: Monitor resources
echo -e "${GREEN}📊 Monitoring resources (Ctrl+C để dừng)...${NC}"
echo ""
echo "Thời gian | Container | CPU % | Memory | Memory % | Status"
echo "---------|-----------|-------|--------|----------|--------"

# Create monitoring function
monitor_resources() {
    while true; do
        TIMESTAMP=$(date '+%H:%M:%S')
        
        # Get stats for all containers
        docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" | tail -n +2 | while read line; do
            CONTAINER=$(echo $line | awk '{print $1}')
            CPU=$(echo $line | awk '{print $2}')
            MEM=$(echo $line | awk '{print $3" "$4}')
            MEM_PERC=$(echo $line | awk '{print $5}')
            
            # Get container status
            STATUS=$(docker inspect --format='{{.State.Status}}' $CONTAINER 2>/dev/null || echo "unknown")
            
            # Color code based on resource usage
            CPU_NUM=$(echo $CPU | sed 's/%//')
            MEM_NUM=$(echo $MEM_PERC | sed 's/%//')
            
            COLOR=$GREEN
            if (( $(echo "$CPU_NUM > 80" | bc -l 2>/dev/null || echo 0) )); then
                COLOR=$RED
            elif (( $(echo "$CPU_NUM > 50" | bc -l 2>/dev/null || echo 0) )); then
                COLOR=$YELLOW
            fi
            
            if (( $(echo "$MEM_NUM > 80" | bc -l 2>/dev/null || echo 0) )); then
                COLOR=$RED
            fi
            
            # Check if container is unhealthy or stopped
            if [ "$STATUS" != "running" ]; then
                COLOR=$RED
                echo -e "${COLOR}$TIMESTAMP | $CONTAINER | $CPU | $MEM | $MEM_PERC | ❌ $STATUS${NC}"
                
                # Show logs if container died
                echo -e "${RED}📋 Last 20 lines of logs:${NC}"
                docker logs --tail 20 $CONTAINER 2>&1 | sed 's/^/  /'
                echo ""
            else
                echo -e "${COLOR}$TIMESTAMP | $CONTAINER | $CPU | $MEM | $MEM_PERC | ✅ $STATUS${NC}"
            fi
        done
        
        echo "---"
        sleep 5
    done
}

# Start monitoring in background
monitor_resources &
MONITOR_PID=$!

# Wait for user interrupt
wait $MONITOR_PID
