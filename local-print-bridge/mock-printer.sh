#!/bin/bash

# Mock Printer Server for Testing
# Listens on port 9100 and accepts print data

PORT=${1:-9100}

echo "Starting Mock Printer Server on port $PORT..."
echo "Press Ctrl+C to stop"
echo ""

# Use netcat to listen on port 9100
while true; do
    echo "Waiting for print job..."
    nc -l $PORT | {
        echo ""
        echo "=========================================="
        echo "📄 Received Print Job at $(date)"
        echo "=========================================="
        
        # Read and display data
        data=$(cat)
        size=${#data}
        
        echo "Data size: $size bytes"
        echo ""
        
        # Show first 200 bytes
        if [ $size -gt 200 ]; then
            echo "First 200 bytes:"
            echo "$data" | head -c 200
            echo ""
            echo "... (truncated)"
        else
            echo "Data:"
            echo "$data"
        fi
        
        echo ""
        echo "=========================================="
        echo "✅ Print job completed"
        echo "=========================================="
        echo ""
    }
done
