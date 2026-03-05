#!/bin/bash

echo "=== Testing Chromedp Optimized Capture ==="
echo ""

cd backend

echo "Running test..."
go run cmd/test-chromedp-optimized/main.go

echo ""
echo "=== Checking output files ==="

if [ -f "test_bill_optimized.bin" ]; then
    SIZE=$(wc -c < test_bill_optimized.bin)
    echo "✅ test_bill_optimized.bin: $SIZE bytes"
else
    echo "❌ test_bill_optimized.bin not found"
fi

if [ -f "test_bill_optimized_preview.png" ]; then
    SIZE=$(wc -c < test_bill_optimized_preview.png)
    echo "✅ test_bill_optimized_preview.png: $SIZE bytes"
    
    # Check image dimensions if 'file' command is available
    if command -v file &> /dev/null; then
        file test_bill_optimized_preview.png
    fi
else
    echo "❌ test_bill_optimized_preview.png not found"
fi

echo ""
echo "=== Test completed ==="
