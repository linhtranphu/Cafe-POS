#!/bin/bash

echo "🧪 Test Print Job - Direct HTTP"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test data
PRINTER_IP="192.168.1.115"

echo "📝 Sending print job to Print Bridge..."
echo "   Printer IP: $PRINTER_IP"
echo "   Port: 9100"
echo ""

# Send directly to Print Bridge HTTP endpoint
RESPONSE=$(curl -s -X POST http://localhost:3001/print \
  -H "Content-Type: application/json" \
  -d '{
    "content": "=== TEST PRINT ===\nWebSocket Test\nTimestamp: '$(date +%Y-%m-%d\ %H:%M:%S)'\n",
    "printerIp": "'$PRINTER_IP'",
    "printerPort": 9100
  }')

echo "Response: $RESPONSE"
echo ""

if echo "$RESPONSE" | grep -q "success"; then
    echo "✅ Print job sent successfully!"
    echo ""
    echo "📋 Check Print Bridge logs:"
    echo "   docker logs local-print-bridge"
    echo "   Or if running locally, check the terminal output"
else
    echo "❌ Failed to send print job"
    echo "   Response: $RESPONSE"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
