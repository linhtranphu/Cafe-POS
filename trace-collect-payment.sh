#!/bin/bash

echo "=== Trace Collect Payment Flow ==="
echo ""

echo "Hướng dẫn:"
echo "1. Chạy script này"
echo "2. Script sẽ theo dõi logs realtime"
echo "3. Bạn collect payment một order"
echo "4. Xem logs có gì"
echo ""
echo "Press Ctrl+C để dừng"
echo ""
echo "Bắt đầu theo dõi logs..."
echo "================================"
echo ""

# Follow logs và filter các dòng quan trọng
docker logs -f cafe-pos-backend 2>&1 | grep --line-buffered -i "payment\|print\|order.*paid\|error\|failed\|creating"
