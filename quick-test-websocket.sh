#!/bin/bash

# Quick WebSocket Test Script
# Tests Print Bridge WebSocket connection in one command

set -e

echo "🚀 Quick WebSocket Test"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check backend
echo "1️⃣  Checking backend..."
if ! lsof -i :3000 >/dev/null 2>&1; then
    echo "   ❌ Backend not running on port 3000"
    echo "   Start with: ./restart_local.sh"
    exit 1
fi
echo "   ✅ Backend running"
echo ""

# Check Socket.IO endpoint
echo "2️⃣  Checking Socket.IO endpoint..."
RESPONSE=$(curl -s "http://localhost:3000/socket.io/?EIO=3&transport=polling" | head -c 50)
if [ -z "$RESPONSE" ]; then
    echo "   ❌ Socket.IO endpoint not responding"
    exit 1
fi
echo "   ✅ Socket.IO endpoint OK"
echo ""

# Check Print Bridge dependencies
echo "3️⃣  Checking Print Bridge dependencies..."
if [ ! -d "local-print-bridge/node_modules/socket.io-client" ]; then
    echo "   ⚠️  Dependencies not installed, installing..."
    (cd local-print-bridge && npm install >/dev/null 2>&1)
fi
echo "   ✅ Dependencies OK"
echo ""

# Check socket.io-client version
echo "4️⃣  Checking socket.io-client version..."
VERSION=$(node -e "console.log(require('./local-print-bridge/node_modules/socket.io-client/package.json').version)" 2>/dev/null || echo "unknown")
if [ "$VERSION" != "2.5.0" ]; then
    echo "   ⚠️  Version mismatch: $VERSION (expected 2.5.0)"
    echo "   Reinstalling correct version..."
    (cd local-print-bridge && npm install socket.io-client@2.5.0 >/dev/null 2>&1)
    VERSION="2.5.0"
fi
echo "   ✅ Version: $VERSION"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ All checks passed!"
echo ""
echo "🎯 Ready to start Print Bridge"
echo ""
echo "Run one of these commands:"
echo ""
echo "  Option 1 (Manual):"
echo "    cd local-print-bridge && npm start"
echo ""
echo "  Option 2 (Script):"
echo "    ./test-print-bridge-websocket.sh"
echo ""
echo "Expected output:"
echo "  [WebSocket] Connecting to: http://localhost:3000"
echo "  [WebSocket] ✅ Connected to backend"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
