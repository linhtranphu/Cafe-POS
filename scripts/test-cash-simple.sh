#!/bin/bash

echo "Testing cash payment update..."
echo ""
echo "Please do manually:"
echo "1. Login as waiter (waiter1/password123)"
echo "2. Start a shift"
echo "3. Create an order"
echo "4. Pay with CASH"
echo "5. Check backend log for DEBUG messages"
echo ""
echo "Watching backend log..."
echo ""

tail -f backend/server.log | grep --line-buffered -E "DEBUG|ERROR|Payment"
