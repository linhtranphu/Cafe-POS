#!/bin/bash

# Stop Local Print Bridge

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "Stopping Local Print Bridge..."
echo ""

# Check if PID file exists
if [ -f .print-bridge.pid ]; then
    PID=$(cat .print-bridge.pid)
    
    # Check if process is running
    if ps -p $PID > /dev/null 2>&1; then
        echo "Killing process $PID..."
        kill $PID
        
        # Wait for process to stop
        sleep 1
        
        # Force kill if still running
        if ps -p $PID > /dev/null 2>&1; then
            echo "Force killing process $PID..."
            kill -9 $PID
        fi
        
        echo -e "${GREEN}✅ Print Bridge stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Process $PID is not running${NC}"
    fi
    
    # Remove PID file
    rm .print-bridge.pid
else
    # Try to find process by port
    if lsof -i :3001 > /dev/null 2>&1; then
        echo "Found process on port 3001:"
        lsof -i :3001
        echo ""
        
        PID=$(lsof -ti :3001)
        echo "Killing process $PID..."
        kill -9 $PID
        
        echo -e "${GREEN}✅ Print Bridge stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Print Bridge is not running${NC}"
    fi
fi

echo ""
echo "Port 3001 status:"
lsof -i :3001 || echo "Port 3001 is free"
