#!/bin/bash

echo "=== Testing with Real Logo from Settings API ==="
echo ""

# Check if backend is running
if ! curl -s http://localhost:3000/api/settings > /dev/null; then
    echo "❌ Backend is not running at http://localhost:3000"
    echo "Please start the backend first:"
    echo "  ./restart_local.sh"
    exit 1
fi

echo "✅ Backend is running"
echo ""

cd backend

echo "Running test with real logo..."
go run cmd/test-real-logo/main.go

echo ""
echo "=== Checking output files ==="

if [ -f "test_real_logo.bin" ]; then
    SIZE=$(wc -c < test_real_logo.bin)
    echo "✅ test_real_logo.bin: $SIZE bytes"
else
    echo "❌ test_real_logo.bin not found"
fi

if [ -f "test_real_logo_preview.png" ]; then
    SIZE=$(wc -c < test_real_logo_preview.png)
    echo "✅ test_real_logo_preview.png: $SIZE bytes"
    
    # Check image dimensions
    if command -v file &> /dev/null; then
        file test_real_logo_preview.png
    fi
    
    echo ""
    echo "Opening preview image..."
    open test_real_logo_preview.png
else
    echo "❌ test_real_logo_preview.png not found"
fi

if [ -f "debug_rendered.html" ]; then
    echo "✅ debug_rendered.html available for inspection"
fi

echo ""
echo "=== Test completed ==="
