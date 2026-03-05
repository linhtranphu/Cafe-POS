#!/bin/bash

# Script để chẩn đoán lỗi Print Bridge không available
# Usage: ./diagnose-print-bridge.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Print Bridge Diagnostic Tool ===${NC}"
echo ""

# Function to check status
check_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
        return 0
    else
        echo -e "${RED}✗ $2${NC}"
        return 1
    fi
}

# Function to print section
print_section() {
    echo ""
    echo -e "${BLUE}=== $1 ===${NC}"
}

# 1. Check Backend Configuration
print_section "1. Backend Configuration"

if [ -f ".env" ]; then
    echo "Checking .env file..."
    
    PRINT_BRIDGE_URL=$(grep "^PRINT_BRIDGE_URL=" .env | cut -d '=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -z "$PRINT_BRIDGE_URL" ]; then
        echo -e "${RED}✗ PRINT_BRIDGE_URL not found in .env${NC}"
        echo "  Add this line to .env:"
        echo "  PRINT_BRIDGE_URL=http://localhost:3001"
    else
        echo -e "${GREEN}✓ PRINT_BRIDGE_URL found: $PRINT_BRIDGE_URL${NC}"
    fi
else
    echo -e "${RED}✗ .env file not found${NC}"
    PRINT_BRIDGE_URL="http://localhost:3001"
fi

# Default to localhost if not set
if [ -z "$PRINT_BRIDGE_URL" ]; then
    PRINT_BRIDGE_URL="http://localhost:3001"
    echo -e "${YELLOW}Using default: $PRINT_BRIDGE_URL${NC}"
fi

# Extract host and port
BRIDGE_HOST=$(echo $PRINT_BRIDGE_URL | sed -e 's|^[^/]*//||' -e 's|/.*$||' -e 's|:.*$||')
BRIDGE_PORT=$(echo $PRINT_BRIDGE_URL | sed -e 's|^[^:]*:||' -e 's|/.*$||')

if [ "$BRIDGE_PORT" = "$BRIDGE_HOST" ]; then
    BRIDGE_PORT="80"
fi

echo "  Host: $BRIDGE_HOST"
echo "  Port: $BRIDGE_PORT"

# 2. Check Print Bridge Service
print_section "2. Print Bridge Service Status"

# Check if running locally
if [ -d "local-print-bridge" ]; then
    echo "Checking local print bridge..."
    
    cd local-print-bridge
    
    # Check Docker container
    if docker ps --format '{{.Names}}' | grep -q "local-print-bridge"; then
        echo -e "${GREEN}✓ Docker container is running${NC}"
        
        # Show container info
        echo ""
        echo "Container details:"
        docker ps --filter "name=local-print-bridge" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
        
        # Check logs for errors
        echo ""
        echo "Recent logs (last 10 lines):"
        docker logs --tail=10 local-print-bridge 2>&1 | sed 's/^/  /'
    else
        echo -e "${RED}✗ Docker container is not running${NC}"
        
        # Check if container exists but stopped
        if docker ps -a --format '{{.Names}}' | grep -q "local-print-bridge"; then
            echo -e "${YELLOW}Container exists but is stopped${NC}"
            echo ""
            echo "To start it:"
            echo "  cd local-print-bridge"
            echo "  docker compose up -d"
        else
            echo -e "${YELLOW}Container does not exist${NC}"
            echo ""
            echo "To create and start it:"
            echo "  cd local-print-bridge"
            echo "  docker compose up -d"
        fi
    fi
    
    cd ..
else
    echo -e "${YELLOW}local-print-bridge directory not found${NC}"
    echo "Print bridge might be running on a different machine"
fi

# 3. Check Network Connectivity
print_section "3. Network Connectivity"

echo "Testing connection to $PRINT_BRIDGE_URL..."

# Ping host
echo ""
echo "Ping test:"
if ping -c 1 -W 2 $BRIDGE_HOST > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Host $BRIDGE_HOST is reachable${NC}"
else
    echo -e "${RED}✗ Cannot reach host $BRIDGE_HOST${NC}"
fi

# Check port
echo ""
echo "Port test:"
if command -v nc > /dev/null 2>&1; then
    if nc -z -w 2 $BRIDGE_HOST $BRIDGE_PORT 2>/dev/null; then
        echo -e "${GREEN}✓ Port $BRIDGE_PORT is open${NC}"
    else
        echo -e "${RED}✗ Port $BRIDGE_PORT is not accessible${NC}"
        echo "  Possible causes:"
        echo "  - Print bridge is not running"
        echo "  - Firewall blocking the port"
        echo "  - Wrong host/port in PRINT_BRIDGE_URL"
    fi
else
    echo -e "${YELLOW}nc (netcat) not available, skipping port test${NC}"
fi

# 4. Check Health Endpoint
print_section "4. Health Check"

echo "Testing $PRINT_BRIDGE_URL/health..."

HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" --connect-timeout 5 "$PRINT_BRIDGE_URL/health" 2>/dev/null || echo "000")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$HEALTH_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    echo "  Response: $RESPONSE_BODY"
elif [ "$HTTP_CODE" = "000" ]; then
    echo -e "${RED}✗ Cannot connect to print bridge${NC}"
    echo "  Error: Connection failed or timeout"
else
    echo -e "${RED}✗ Health check failed${NC}"
    echo "  HTTP Code: $HTTP_CODE"
    echo "  Response: $RESPONSE_BODY"
fi

# 5. Check Backend Service
print_section "5. Backend Service"

# Check if backend is running
if pgrep -f "main" > /dev/null || docker ps --format '{{.Names}}' | grep -q "backend"; then
    echo -e "${GREEN}✓ Backend service is running${NC}"
    
    # Check backend logs for print bridge errors
    echo ""
    echo "Checking backend logs for print bridge errors..."
    
    if docker ps --format '{{.Names}}' | grep -q "backend"; then
        echo "Recent print bridge related logs:"
        docker logs --tail=50 backend 2>&1 | grep -i "print\|bridge" | tail -10 | sed 's/^/  /' || echo "  No print bridge related logs found"
    fi
else
    echo -e "${YELLOW}Backend service status unknown${NC}"
fi

# 6. Test Print Bridge Endpoints
print_section "6. Test Print Bridge Endpoints"

# Test /health
echo "Testing /health endpoint..."
curl -s -w "\nHTTP Code: %{http_code}\n" --connect-timeout 5 "$PRINT_BRIDGE_URL/health" 2>/dev/null || echo -e "${RED}Failed to connect${NC}"

echo ""

# Test /test-connection
echo "Testing /test-connection endpoint..."
TEST_RESPONSE=$(curl -s -w "\nHTTP Code: %{http_code}\n" --connect-timeout 5 \
    -X POST "$PRINT_BRIDGE_URL/test-connection" \
    -H "Content-Type: application/json" \
    -d '{"printerIP": "192.168.1.100", "printerPort": 9100}' 2>/dev/null || echo -e "${RED}Failed to connect${NC}")

echo "$TEST_RESPONSE"

# 7. Summary and Recommendations
print_section "7. Summary and Recommendations"

echo ""
echo "Common issues and solutions:"
echo ""

echo "1. Print Bridge not running:"
echo "   cd local-print-bridge"
echo "   docker compose up -d"
echo ""

echo "2. Wrong PRINT_BRIDGE_URL in backend .env:"
echo "   Edit .env and set:"
echo "   PRINT_BRIDGE_URL=http://localhost:3001  (if same machine)"
echo "   PRINT_BRIDGE_URL=http://192.168.1.X:3001  (if different machine)"
echo ""

echo "3. Firewall blocking port 3001:"
echo "   sudo ufw allow 3001/tcp"
echo ""

echo "4. Print Bridge crashed:"
echo "   cd local-print-bridge"
echo "   docker compose logs"
echo "   docker compose restart"
echo ""

echo "5. Backend cannot reach Print Bridge (different machines):"
echo "   - Check network connectivity"
echo "   - Verify IP address is correct"
echo "   - Check firewall on Print Bridge machine"
echo ""

# 8. Quick Fix Commands
print_section "8. Quick Fix Commands"

echo ""
echo "Try these commands to fix common issues:"
echo ""

echo "# Restart Print Bridge"
echo "cd local-print-bridge && docker compose restart && cd .."
echo ""

echo "# View Print Bridge logs"
echo "cd local-print-bridge && docker compose logs -f"
echo ""

echo "# Restart Backend"
echo "docker compose restart backend"
echo ""

echo "# Test connection manually"
echo "curl $PRINT_BRIDGE_URL/health"
echo ""

# 9. Environment Info
print_section "9. Environment Information"

echo "Docker version:"
docker --version 2>/dev/null || echo "Docker not installed"

echo ""
echo "Docker Compose version:"
docker compose version 2>/dev/null || echo "Docker Compose not installed"

echo ""
echo "Network interfaces:"
ip addr show 2>/dev/null | grep "inet " | grep -v "127.0.0.1" | awk '{print "  " $2}' || \
    ifconfig 2>/dev/null | grep "inet " | grep -v "127.0.0.1" | awk '{print "  " $2}'

echo ""
echo -e "${BLUE}=== Diagnostic Complete ===${NC}"
echo ""
