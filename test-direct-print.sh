#!/bin/bash

echo "=== Test Direct Print to Printer ==="
echo ""

# Get printer IP from user
read -p "Enter printer IP address (e.g., 192.168.1.100): " PRINTER_IP

if [ -z "$PRINTER_IP" ]; then
    echo "❌ Printer IP is required"
    exit 1
fi

echo ""
echo "Testing with printer: $PRINTER_IP:9100"
echo ""

cd backend

# Test 1: Simple text
echo "=== Test 1: Simple Text ==="
echo "Sending simple text directly to printer..."

# Create simple ESC/POS command
printf "\x1b\x40Hello from Cafe POS!\nDirect print test\n\n\x1b\x64\x03\x1d\x56\x00" | nc -w 3 $PRINTER_IP 9100

if [ $? -eq 0 ]; then
    echo "✅ Test 1 sent"
else
    echo "❌ Test 1 failed - check printer IP and port"
    exit 1
fi

sleep 2

# Test 2: Simple ESC/POS file
echo ""
echo "=== Test 2: Simple ESC/POS (GS v 0) ==="

if [ ! -f "test_simple_escpos.bin" ]; then
    echo "Generating test file..."
    go run cmd/test-simple-escpos/main.go > /dev/null 2>&1
fi

if [ -f "test_simple_escpos.bin" ]; then
    SIZE=$(wc -c < test_simple_escpos.bin)
    echo "File size: $SIZE bytes"
    echo "Sending to printer..."
    
    cat test_simple_escpos.bin | nc -w 3 $PRINTER_IP 9100
    
    if [ $? -eq 0 ]; then
        echo "✅ Test 2 sent"
    else
        echo "❌ Test 2 failed"
    fi
else
    echo "❌ test_simple_escpos.bin not found"
fi

sleep 2

# Test 3: Full bill with logo
echo ""
echo "=== Test 3: Full Bill (GS 8 L) ==="

if [ ! -f "test_uploaded_logo.bin" ]; then
    echo "Generating test file..."
    go run cmd/test-uploaded-logo/main.go > /dev/null 2>&1
fi

if [ -f "test_uploaded_logo.bin" ]; then
    SIZE=$(wc -c < test_uploaded_logo.bin)
    echo "File size: $SIZE bytes"
    echo "Sending to printer..."
    
    cat test_uploaded_logo.bin | nc -w 3 $PRINTER_IP 9100
    
    if [ $? -eq 0 ]; then
        echo "✅ Test 3 sent"
    else
        echo "❌ Test 3 failed"
    fi
else
    echo "❌ test_uploaded_logo.bin not found"
fi

cd ..

echo ""
echo "=== Test completed ==="
echo ""
echo "Check your printer output:"
echo "  Test 1: Simple text message"
echo "  Test 2: Black rectangle (GS v 0 command)"
echo "  Test 3: Full bill with logo (GS 8 L command)"
echo ""
echo "If Test 1 works but Test 2/3 don't, there may be an issue with binary data."
echo "If Test 2 works but Test 3 doesn't, the printer doesn't support GS 8 L."
