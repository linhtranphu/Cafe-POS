#!/bin/bash

echo "🔄 Closing All Open Waiter Shifts (Auto Mode)"
echo "=============================================="
echo ""
echo "⚠️  This will automatically close ALL open waiter shifts!"
echo ""

# Change to backend directory and run
(cd backend && go run cmd/close-all-waiter-shifts/main.go)
