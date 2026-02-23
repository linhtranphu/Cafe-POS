#!/bin/bash

echo "🔄 Closing All Open Waiter Shifts"
echo "=================================="
echo ""
echo "⚠️  WARNING: This will close ALL open waiter shifts!"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
  echo "❌ Operation cancelled"
  exit 0
fi

echo ""
echo "Running script..."
echo ""

cd backend
go run cmd/close-all-waiter-shifts/main.go
