#!/bin/bash

echo "🧪 Testing Payment Shift Update (Direct DB Test)"
echo "================================================"
echo ""

cd backend
go run cmd/test-payment-shift-update/main.go
