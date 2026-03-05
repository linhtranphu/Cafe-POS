#!/bin/bash

# Debug Print Flow
# Tests the entire print flow from frontend to printer

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

PRINTER_IP="${1:-192.168.1.115}"
PRINTER_PORT="${2:-9100}"

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║          Debug Print Flow                              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

echo "Testing printer: ${PRINTER_IP}:${PRINTER_PORT}"
echo ""

# Step 1: Check local network
echo -e "${YELLOW}Step 1: Check local network${NC}"
echo "Your IP:"
ifconfig | grep "inet " | grep -v 127.0.0.1
echo ""

# Step 2: Ping printer
echo -e "${YELLOW}Step 2: Ping printer${NC}"
ping -c 2 -W 2 ${PRINTER_IP}
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Printer is reachable${NC}"
else
    echo -e "${RED}❌ Cannot ping printer${NC}"
    echo "Possible issues:"
    echo "  - Printer is turned off"
    echo "  - Wrong IP address"
    echo "  - Printer on different network"
    echo "  - Firewall blocking ICMP"
fi
echo ""

# Step 3: Test printer port
echo -e "${YELLOW}Step 3: Test printer port ${PRINTER_PORT}${NC}"
timeout 3 nc -zv ${PRINTER_IP} ${PRINTER_PORT} 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Printer port is open${NC}"
else
    echo -e "${RED}❌ Cannot connect to printer port${NC}"
    echo "Possible issues:"
    echo "  - Printer service not running"
    echo "  - Wrong port (try 9100, 515, or check printer manual)"
    echo "  - Firewall blocking port"
fi
echo ""

# Step 4: Check print bridge
echo -e "${YELLOW}Step 4: Check local print bridge${NC}"
if lsof -i :3001 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Print bridge is running on port 3001${NC}"
    
    # Test health
    response=$(curl -s http://localhost:3001/health)
    echo "Health check: $response"
else
    echo -e "${RED}❌ Print bridge is not running${NC}"
    echo "Start it with: cd local-print-bridge && ./print-bridge"
fi
echo ""

# Step 5: Check Cloudflare tunnel
echo -e "${YELLOW}Step 5: Check Cloudflare tunnel${NC}"
if pgrep -f "cloudflared tunnel" > /dev/null; then
    echo -e "${GREEN}✅ Cloudflare tunnel is running${NC}"
    
    # Test public URL
    response=$(curl -s https://print.tacafe.store/health 2>&1)
    if echo "$response" | grep -q "ok"; then
        echo -e "${GREEN}✅ Public URL is accessible${NC}"
        echo "URL: https://print.tacafe.store"
    else
        echo -e "${RED}❌ Public URL not accessible${NC}"
        echo "Response: $response"
    fi
else
    echo -e "${RED}❌ Cloudflare tunnel is not running${NC}"
    echo "Start it with: cd local-print-bridge && ./start-tunnel.sh"
fi
echo ""

# Step 6: Test print bridge connection to printer
echo -e "${YELLOW}Step 6: Test print bridge → printer connection${NC}"
test_response=$(curl -s -X POST http://localhost:3001/test-connection \
    -H "Content-Type: application/json" \
    -d "{\"printerIP\": \"${PRINTER_IP}\", \"printerPort\": ${PRINTER_PORT}}")

if echo "$test_response" | grep -q "success.*true"; then
    echo -e "${GREEN}✅ Print bridge can connect to printer${NC}"
else
    echo -e "${RED}❌ Print bridge cannot connect to printer${NC}"
    echo "Response: $test_response"
fi
echo ""

# Step 7: Scan for printers on network
echo -e "${YELLOW}Step 7: Scan for printers on local network${NC}"
echo "Scanning 192.168.1.1-254 on port 9100..."
echo "(This may take a minute)"
echo ""

for i in {100..120}; do
    ip="192.168.1.$i"
    timeout 0.5 nc -zv $ip 9100 2>&1 | grep -q "succeeded" && echo "Found printer at: $ip:9100"
done

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Summary                             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "If printer is not reachable:"
echo "1. Check printer is turned on"
echo "2. Check printer IP address (print network config page)"
echo "3. Ensure printer and computer are on same network"
echo "4. Try connecting from printer's web interface"
echo "5. Check printer manual for correct port (usually 9100)"
echo ""
echo "If print bridge is not running:"
echo "  cd local-print-bridge && ./print-bridge"
echo ""
echo "If tunnel is not running:"
echo "  cd local-print-bridge && ./start-tunnel.sh"
echo ""
