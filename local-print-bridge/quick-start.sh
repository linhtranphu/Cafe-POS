#!/bin/bash

# Quick Start Script for Local Print Bridge

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   🖨️  Local Print Bridge - Quick Start               ║${NC}"
echo -e "${BLUE}║              (Go + chromedp version)                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ .env file created${NC}"
    echo ""
    echo -e "${YELLOW}Please edit .env and configure your printer IPs:${NC}"
    echo "  nano .env"
    echo ""
    read -p "Press Enter to continue..."
fi

# Check if binary exists
if [ ! -f print-bridge ]; then
    echo -e "${YELLOW}⚠️  Binary not found. Building...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed!${NC}"
        echo "Please install Go from: https://golang.org/dl/"
        exit 1
    fi
    
    echo "Building print-bridge..."
    go build -o print-bridge main.go
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Build failed!${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Build successful${NC}"
    echo ""
fi

# Check if port 3001 is already in use
if lsof -i :3001 > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Port 3001 is already in use!${NC}"
    echo ""
    echo "Process using port 3001:"
    lsof -i :3001
    echo ""
    read -p "Kill existing process? (y/n): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        PID=$(lsof -ti :3001)
        kill -9 $PID
        echo -e "${GREEN}✅ Process killed${NC}"
        sleep 1
    else
        echo "Please stop the existing process first."
        exit 1
    fi
fi

# Start print bridge
echo -e "${GREEN}🚀 Starting Local Print Bridge...${NC}"
echo ""

# Run in background
./print-bridge &
PID=$!

echo "Process ID: $PID"
echo ""

# Wait for service to start
echo "Waiting for service to start..."
sleep 2

# Test health check
echo ""
echo -e "${YELLOW}Testing health check...${NC}"
response=$(curl -s http://localhost:3001/health)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Print Bridge is running!${NC}"
    echo ""
    echo "Response:"
    echo "$response" | jq . 2>/dev/null || echo "$response"
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                Service Information                     ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "  URL: http://localhost:3001"
    echo "  PID: $PID"
    echo ""
    echo -e "${YELLOW}Available Endpoints:${NC}"
    echo "  GET  /health              - Health check"
    echo "  GET  /status              - Service status"
    echo "  POST /render-and-print    - Render HTML and print"
    echo "  POST /print               - Direct print ESC/POS"
    echo "  POST /test-connection     - Test printer connection"
    echo ""
    echo -e "${YELLOW}Test Commands:${NC}"
    echo "  # Health check"
    echo "  curl http://localhost:3001/health"
    echo ""
    echo "  # Status"
    echo "  curl http://localhost:3001/status"
    echo ""
    echo "  # Test printer connection"
    echo "  curl -X POST http://localhost:3001/test-connection \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"printerIP\": \"192.168.1.100\"}'"
    echo ""
    echo -e "${YELLOW}Stop Service:${NC}"
    echo "  kill $PID"
    echo ""
    echo -e "${GREEN}✅ Print Bridge started successfully!${NC}"
    echo ""
    
    # Save PID to file
    echo $PID > .print-bridge.pid
    echo "PID saved to .print-bridge.pid"
    echo ""
else
    echo -e "${RED}❌ Failed to start Print Bridge!${NC}"
    echo ""
    echo "Check logs for errors:"
    echo "  tail -f logs/print-bridge.log"
    
    # Kill the process
    kill $PID 2>/dev/null
    exit 1
fi
