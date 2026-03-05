#!/bin/bash

echo "=== Checking Backend Status on EC2 ==="
echo ""

# Check if backend is running
echo "1. Backend container status:"
ssh -i ~/.ssh/tacafe-ec2.pem ubuntu@tacafe.store "docker ps | grep backend"
echo ""

# Check backend logs for print bridge
echo "2. Backend logs (print bridge related):"
ssh -i ~/.ssh/tacafe-ec2.pem ubuntu@tacafe.store "docker logs cafe-pos-backend 2>&1 | grep -i 'print bridge\|chromedp\|html template' | tail -20"
echo ""

# Check environment variables
echo "3. Backend environment (PRINT_BRIDGE_URL):"
ssh -i ~/.ssh/tacafe-ec2.pem ubuntu@tacafe.store "docker exec cafe-pos-backend env | grep PRINT_BRIDGE"
echo ""

# Check if .env has PRINT_BRIDGE_URL
echo "4. .env file (PRINT_BRIDGE_URL):"
ssh -i ~/.ssh/tacafe-ec2.pem ubuntu@tacafe.store "cat ~/cafe-pos/.env | grep PRINT_BRIDGE"
echo ""

# Check Socket.IO status
echo "5. Socket.IO server status:"
ssh -i ~/.ssh/tacafe-ec2.pem ubuntu@tacafe.store "docker logs cafe-pos-backend 2>&1 | grep -i 'socket.io' | tail -10"
echo ""

echo "=== Diagnosis ==="
echo ""
echo "If PRINT_BRIDGE_URL is missing:"
echo "  → HTML template routes will NOT be registered (404 error)"
echo "  → Need to add PRINT_BRIDGE_URL to .env and restart backend"
echo ""
echo "If Socket.IO not started:"
echo "  → WebSocket connections will fail"
echo "  → Need to check backend code and restart"
