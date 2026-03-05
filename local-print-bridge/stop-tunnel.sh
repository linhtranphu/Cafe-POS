#!/bin/bash

# Stop Cloudflare Tunnel

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "Stopping Cloudflare Tunnel..."
echo ""

# Check if PID file exists
if [ -f .tunnel.pid ]; then
    PID=$(cat .tunnel.pid)
    
    # Check if process is running
    if ps -p $PID > /dev/null 2>&1; then
        echo "Killing tunnel process $PID..."
        kill $PID
        
        # Wait for process to stop
        sleep 2
        
        # Force kill if still running
        if ps -p $PID > /dev/null 2>&1; then
            echo "Force killing process $PID..."
            kill -9 $PID
        fi
        
        echo -e "${GREEN}✅ Cloudflare Tunnel stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Process $PID is not running${NC}"
    fi
    
    # Remove PID file
    rm .tunnel.pid
else
    # Try to find process by name
    if pgrep -f "cloudflared tunnel run" > /dev/null; then
        echo "Found tunnel process:"
        ps aux | grep "cloudflared tunnel run" | grep -v grep
        echo ""
        
        PID=$(pgrep -f "cloudflared tunnel run")
        echo "Killing process $PID..."
        pkill -f "cloudflared tunnel run"
        
        echo -e "${GREEN}✅ Cloudflare Tunnel stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  Cloudflare Tunnel is not running${NC}"
    fi
fi

echo ""
echo "Tunnel status:"
pgrep -f "cloudflared tunnel run" > /dev/null && echo "Still running" || echo "Stopped"
