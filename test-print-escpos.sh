#!/bin/bash

echo "=== Test Print ESC/POS Files ==="
echo ""

# Check if print bridge is running
if ! curl -s http://localhost:3001/health > /dev/null 2>&1; then
    echo "❌ Print bridge is not running at http://localhost:3001"
    echo "Please start it with: docker start local-print-bridge"
    exit 1
fi

echo "✅ Print bridge is running"
echo ""

# Get printer IP from user
read -p "Enter printer IP address (e.g., 192.168.1.100): " PRINTER_IP

if [ -z "$PRINTER_IP" ]; then
    echo "❌ Printer IP is required"
    exit 1
fi

echo ""
echo "Testing with printer: $PRINTER_IP"
echo ""

cd backend

# Test 1: Simple ESC/POS (GS v 0 - Raster bit image)
echo "=== Test 1: Simple ESC/POS (GS v 0 command) ==="
if [ ! -f "test_simple_escpos.bin" ]; then
    echo "Generating test file..."
    go run cmd/test-simple-escpos/main.go
fi

if [ -f "test_simple_escpos.bin" ]; then
    SIZE=$(wc -c < test_simple_escpos.bin)
    echo "File size: $SIZE bytes"
    
    # Convert to base64
    BASE64_DATA=$(base64 < test_simple_escpos.bin | tr -d '\n')
    
    echo "Sending to printer via print bridge..."
    RESPONSE=$(curl -s -X POST http://localhost:3001/print \
        -H "Content-Type: application/json" \
        -d "{
            \"printer_ip\": \"$PRINTER_IP\",
            \"data\": \"$BASE64_DATA\",
            \"encoding\": \"base64\"
        }")
    
    echo "Response: $RESPONSE"
    
    if echo "$RESPONSE" | grep -q "success"; then
        echo "✅ Test 1 sent successfully"
    else
        echo "❌ Test 1 failed"
    fi
else
    echo "❌ test_simple_escpos.bin not found"
fi

echo ""
echo "=== Test 2: ESC/POS Library (GS 8 L command) ==="

# Generate if not exists
if [ ! -f "test_uploaded_logo.bin" ]; then
    echo "Generating test file with logo..."
    go run cmd/test-uploaded-logo/main.go > /dev/null 2>&1
fi

if [ -f "test_uploaded_logo.bin" ]; then
    SIZE=$(wc -c < test_uploaded_logo.bin)
    echo "File size: $SIZE bytes"
    
    # Convert to base64
    BASE64_DATA=$(base64 < test_uploaded_logo.bin | tr -d '\n')
    
    echo "Sending to printer via print bridge..."
    RESPONSE=$(curl -s -X POST http://localhost:3001/print \
        -H "Content-Type: application/json" \
        -d "{
            \"printer_ip\": \"$PRINTER_IP\",
            \"data\": \"$BASE64_DATA\",
            \"encoding\": \"base64\"
        }")
    
    echo "Response: $RESPONSE"
    
    if echo "$RESPONSE" | grep -q "success"; then
        echo "✅ Test 2 sent successfully"
    else
        echo "❌ Test 2 failed"
    fi
else
    echo "❌ test_uploaded_logo.bin not found"
fi

cd ..

echo ""
echo "=== Test completed ==="
echo ""
echo "Check your printer to see which test printed correctly:"
echo "  Test 1: Simple black rectangle (GS v 0 command)"
echo "  Test 2: Full bill with logo (GS 8 L command)"
echo ""
echo "If Test 1 works but Test 2 doesn't, the printer doesn't support GS 8 L"
echo "and we need to switch to GS v 0 (raster bit image) command."
