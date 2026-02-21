#!/bin/bash

echo "=== Testing WebSocket Connection ==="
echo ""

# Test 1: Check if backend is running
echo "1. Checking backend status..."
if curl -s http://localhost:3000/api/state-machines > /dev/null; then
    echo "✅ Backend is running on port 3000"
else
    echo "❌ Backend is not responding"
    exit 1
fi

echo ""

# Test 2: Test Socket.IO HTTP polling handshake (EIO=4)
echo "2. Testing Socket.IO v4 HTTP polling handshake..."
HANDSHAKE_RESPONSE=$(curl -s "http://localhost:3000/socket.io/?EIO=4&transport=polling")
echo "Response: $HANDSHAKE_RESPONSE"

if echo "$HANDSHAKE_RESPONSE" | grep -q "sid"; then
    echo "✅ Socket.IO handshake successful"
else
    echo "❌ Socket.IO handshake failed"
    exit 1
fi

echo ""

# Test 3: Extract session ID
echo "3. Extracting session ID..."
SID=$(echo "$HANDSHAKE_RESPONSE" | grep -o '"sid":"[^"]*"' | cut -d'"' -f4)
echo "Session ID: $SID"

echo ""
echo "=== WebSocket Connection Test Complete ==="
echo ""
echo "Next steps:"
echo "1. Open browser console at http://localhost:5173"
echo "2. Check for WebSocket connection logs"
echo "3. Look for '[WebSocket] Connected' message"
