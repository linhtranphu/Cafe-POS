#!/bin/bash

echo "🧪 Testing Print Bridge WebSocket Connection"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if backend is running
if ! lsof -i :3000 >/dev/null 2>&1; then
    echo "❌ Backend is not running on port 3000"
    echo "   Start it with: ./restart_local.sh"
    exit 1
fi

echo "✅ Backend is running on port 3000"
echo ""

# Check if Print Bridge dependencies are installed
if [ ! -d "local-print-bridge/node_modules" ]; then
    echo "⚠️  Print Bridge dependencies not installed"
    echo "   Installing..."
    cd local-print-bridge && npm install && cd ..
fi

echo "📝 Starting Print Bridge with WebSocket..."
echo "   Watch for connection messages:"
echo "   - [WebSocket] Connecting to: http://localhost:3000"
echo "   - [WebSocket] ✅ Connected to backend"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Start Print Bridge
cd local-print-bridge
npm start
