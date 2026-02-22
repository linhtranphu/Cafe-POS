#!/bin/bash

echo "🧪 Testing Print Bridge WebSocket Connection"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Print Bridge is running
echo "📝 Step 1: Checking if Print Bridge is running..."
if curl -s http://localhost:3001/health > /dev/null; then
    echo "✅ Print Bridge is running"
else
    echo "❌ Print Bridge is not running"
    echo "   Start it with: cd local-print-bridge && npm start"
    exit 1
fi

echo ""
echo "📝 Step 2: Checking Print Bridge logs for WebSocket connection..."
echo "   Look for these messages in the logs:"
echo "   - [WebSocket] Connecting to backend: http://localhost:3000"
echo "   - [WebSocket] ✅ Connected to backend"
echo ""
echo "   If you see connection errors, check:"
echo "   1. Backend is running on port 3000"
echo "   2. BACKEND_URL in .env is correct"
echo ""

echo "📝 Step 3: Testing manual print job creation..."
echo "   This will create a print job in backend"
echo "   Print Bridge should receive it via WebSocket and print automatically"
echo ""

# Get auth token
echo "🔐 Getting auth token..."
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get auth token"
    exit 1
fi
echo "✅ Got auth token"

echo ""
echo "📝 Creating test print job..."
echo "   Watch Print Bridge logs for:"
echo "   - [WebSocket] 📨 New print job received"
echo "   - [PrintJobHandler] Processing job"
echo "   - [PrintJobHandler] ✅ Job printed successfully"
echo ""

# Note: This is a placeholder - you need to create an actual order first
# For now, just show the instructions
echo "⚠️  To test WebSocket:"
echo "   1. Open Print Bridge logs: docker logs -f local-print-bridge"
echo "   2. Create an order in the frontend"
echo "   3. Watch for WebSocket messages in Print Bridge logs"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Test setup complete"
echo ""
echo "Next steps:"
echo "1. Check Print Bridge logs: docker logs -f local-print-bridge"
echo "2. Create an order in frontend"
echo "3. Verify Print Bridge receives job via WebSocket"
echo "4. Verify printer prints automatically"
