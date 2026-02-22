#!/bin/bash

echo "🧪 Testing Print Job WebSocket Broadcast"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "📝 Step 1: Check Print Bridge WebSocket connection..."
echo "   Watch Print Bridge logs in another terminal:"
echo "   tail -f local-print-bridge logs or check process output"
echo ""

echo "📝 Step 2: This test requires:"
echo "   - Backend running on port 3000 ✅"
echo "   - Print Bridge running with WebSocket ✅"
echo "   - MongoDB running ✅"
echo ""

echo "📝 Step 3: To test, you need to:"
echo "   1. Open frontend: http://localhost:5173"
echo "   2. Login with admin/admin123"
echo "   3. Create an order"
echo "   4. Watch Print Bridge logs for:"
echo "      [WebSocket] 📨 New print job received"
echo "      [PrintJobHandler] Processing job"
echo ""

echo "Or manually trigger a print job broadcast from backend..."
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Instructions complete"
echo ""
echo "Next: Create an order in frontend and watch logs!"
