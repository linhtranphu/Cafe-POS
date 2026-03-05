#!/bin/bash

# Start Cloudflare Tunnel for Print Bridge

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Starting Cloudflare Tunnel for Print Bridge         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if cloudflared is installed
if ! command -v cloudflared &> /dev/null; then
    echo -e "${RED}❌ cloudflared is not installed!${NC}"
    echo "Please run: brew install cloudflared"
    exit 1
fi

# Check if config exists
if [ ! -f ~/.cloudflared/config.yml ]; then
    echo -e "${RED}❌ Cloudflare tunnel config not found!${NC}"
    echo "Please run: ./cloudflare-tunnel-setup.sh"
    exit 1
fi

# Get tunnel name from config
TUNNEL_ID=$(grep "^tunnel:" ~/.cloudflared/config.yml | awk '{print $2}')
TUNNEL_NAME=$(cloudflared tunnel list | grep "$TUNNEL_ID" | awk '{print $2}')

echo "Tunnel ID: $TUNNEL_ID"
echo "Tunnel Name: $TUNNEL_NAME"
echo ""

# Check if tunnel is already running
if pgrep -f "cloudflared tunnel run" > /dev/null; then
    echo -e "${YELLOW}⚠️  Cloudflare tunnel is already running!${NC}"
    echo ""
    echo "Process:"
    ps aux | grep "cloudflared tunnel run" | grep -v grep
    echo ""
    read -p "Kill existing process and restart? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        pkill -f "cloudflared tunnel run"
        sleep 2
        echo -e "${GREEN}✅ Existing process killed${NC}"
    else
        exit 0
    fi
fi

# Check if print bridge is running
if ! lsof -i :3001 > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Print Bridge is not running on port 3001!${NC}"
    echo ""
    read -p "Start Print Bridge first? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Starting Print Bridge..."
        ./quick-start.sh &
        sleep 3
    else
        echo -e "${RED}Warning: Tunnel will start but Print Bridge is not running${NC}"
    fi
fi

# Start tunnel in background
echo -e "${GREEN}🚀 Starting Cloudflare Tunnel...${NC}"
echo ""

# Run tunnel in background and save PID
nohup cloudflared tunnel run > tunnel.log 2>&1 &
TUNNEL_PID=$!

echo "Tunnel PID: $TUNNEL_PID"
echo $TUNNEL_PID > .tunnel.pid

# Wait for tunnel to start
echo "Waiting for tunnel to connect..."
sleep 5

# Check if tunnel is running
if ps -p $TUNNEL_PID > /dev/null; then
    echo -e "${GREEN}✅ Cloudflare Tunnel started successfully!${NC}"
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                Tunnel Information                      ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "  Public URL: https://print.tacafe.store"
    echo "  Local Service: http://localhost:3001"
    echo "  Tunnel PID: $TUNNEL_PID"
    echo "  Log File: tunnel.log"
    echo ""
    echo -e "${YELLOW}Test Commands:${NC}"
    echo "  # Test from internet"
    echo "  curl https://print.tacafe.store/health"
    echo ""
    echo "  # View logs"
    echo "  tail -f tunnel.log"
    echo ""
    echo "  # Check tunnel status"
    echo "  cloudflared tunnel info $TUNNEL_NAME"
    echo ""
    echo "  # Stop tunnel"
    echo "  ./stop-tunnel.sh"
    echo ""
    echo -e "${GREEN}✅ Tunnel is ready!${NC}"
else
    echo -e "${RED}❌ Failed to start tunnel!${NC}"
    echo ""
    echo "Check logs:"
    cat tunnel.log
    exit 1
fi
